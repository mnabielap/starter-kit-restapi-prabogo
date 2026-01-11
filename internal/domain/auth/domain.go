package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/palantir/stacktrace"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
	"prabogo/utils/jwt"
	"prabogo/utils/password"
)

type AuthDomain interface {
	Login(ctx context.Context, email, pass string) (*model.User, map[string]interface{}, error)
	Register(ctx context.Context, input model.UserInput) (*model.User, map[string]interface{}, error)
	RefreshToken(ctx context.Context, refreshToken string) (map[string]interface{}, error)
	Logout(ctx context.Context, refreshToken string) error
	
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type authDomain struct {
	db    outbound_port.DatabasePort
	email outbound_port.EmailPort
}

func NewAuthDomain(db outbound_port.DatabasePort, email outbound_port.EmailPort) AuthDomain {
	return &authDomain{
		db:    db,
		email: email,
	}
}

func (d *authDomain) Login(ctx context.Context, email, pass string) (*model.User, map[string]interface{}, error) {
	user, err := d.db.User().FindByEmail(email)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "db error")
	}
	if user == nil || !password.CheckPassword(pass, user.Password) {
		return nil, nil, stacktrace.NewError("incorrect email or password")
	}

	tokens, err := d.generateAndSaveTokens(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (d *authDomain) Register(ctx context.Context, input model.UserInput) (*model.User, map[string]interface{}, error) {
	repo := d.db.User()
	exists, err := repo.ExistsByEmail(input.Email)
	if exists {
		return nil, nil, stacktrace.NewError("email already taken")
	}

	hashed, err := password.HashPassword(input.Password)
	if err != nil {
		return nil, nil, err
	}

	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashed,
		Role:     "user",
	}
	model.UserPrepare(user)

	if err := repo.Create(user); err != nil {
		return nil, nil, err
	}

	tokens, err := d.generateAndSaveTokens(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (d *authDomain) generateAndSaveTokens(userID string) (map[string]interface{}, error) {
	accessToken, refreshToken, accessExp, refreshExp, err := jwt.GenerateAuthTokens(userID)
	if err != nil {
		return nil, err
	}

	// Save Refresh Token
	err = d.db.Token().Create(&model.Token{
		Token:   refreshToken,
		UserID:  userID,
		Type:    model.TokenTypeRefresh,
		Expires: refreshExp,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"access": map[string]interface{}{
			"token":   accessToken,
			"expires": accessExp,
		},
		"refresh": map[string]interface{}{
			"token":   refreshToken,
			"expires": refreshExp,
		},
	}, nil
}

func (d *authDomain) Logout(ctx context.Context, refreshToken string) error {
	tokenRepo := d.db.Token()
	token, err := tokenRepo.FindByToken(refreshToken, model.TokenTypeRefresh)
	if err != nil || token == nil {
		return stacktrace.NewError("token not found")
	}
	return tokenRepo.Delete(token.ID)
}

func (d *authDomain) RefreshToken(ctx context.Context, refreshToken string) (map[string]interface{}, error) {
	tokenRepo := d.db.Token()
	token, err := tokenRepo.FindByToken(refreshToken, model.TokenTypeRefresh)
	if err != nil || token == nil {
		return nil, stacktrace.NewError("please authenticate")
	}

	// Verify JWT validity
	payload, err := jwt.ValidateLocalToken(refreshToken)
	if err != nil {
		return nil, stacktrace.NewError("invalid token")
	}

	// Delete old token
	tokenRepo.Delete(token.ID)

	// Generate new pair
	return d.generateAndSaveTokens(payload.Sub)
}

func (d *authDomain) ForgotPassword(ctx context.Context, email string) error {
	user, err := d.db.User().FindByEmail(email)
	if err != nil || user == nil {
		return nil // Fail silently
	}

	token, exp, err := jwt.GenerateToken(user.ID, 10*time.Minute, model.TokenTypeResetPassword, "") // Secret handled inside
	if err != nil {
		return err
	}

	err = d.db.Token().Create(&model.Token{
		Token:   token,
		UserID:  user.ID,
		Type:    model.TokenTypeResetPassword,
		Expires: exp,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)
	body := fmt.Sprintf("Click here to reset password: %s", resetURL)
	return d.email.SendEmail(email, "Reset Password", body)
}

func (d *authDomain) ResetPassword(ctx context.Context, tokenStr, newPassword string) error {
	tokenRepo := d.db.Token()
	token, err := tokenRepo.FindByToken(tokenStr, model.TokenTypeResetPassword)
	if err != nil || token == nil {
		return stacktrace.NewError("invalid or expired token")
	}

	// Verify JWT
	_, err = jwt.ValidateLocalToken(tokenStr)
	if err != nil {
		return err
	}

	hashed, _ := password.HashPassword(newPassword)
	user, _ := d.db.User().FindByID(token.UserID)
	if user == nil {
		return stacktrace.NewError("user not found")
	}

	user.Password = hashed
	d.db.User().Update(user)
	
	// Consume token
	tokenRepo.DeleteByUserIDAndType(user.ID, model.TokenTypeResetPassword)
	return nil
}