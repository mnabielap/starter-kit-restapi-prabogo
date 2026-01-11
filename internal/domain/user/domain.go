package user

import (
	"context"
	"time"

	"github.com/palantir/stacktrace"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
	"prabogo/utils/password"
)

type UserDomain interface {
	Create(ctx context.Context, input model.UserInput) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetAll(ctx context.Context, filters model.UserFilter, page, limit int, sort string) ([]model.User, int64, error)
	Update(ctx context.Context, id string, input model.UserInput) (*model.User, error)
	Delete(ctx context.Context, id string) error
}

type userDomain struct {
	db outbound_port.DatabasePort
}

func NewUserDomain(db outbound_port.DatabasePort) UserDomain {
	return &userDomain{db: db}
}

func (d *userDomain) Create(ctx context.Context, input model.UserInput) (*model.User, error) {
	repo := d.db.User()

	exists, err := repo.ExistsByEmail(input.Email)
	if err != nil {
		return nil, stacktrace.Propagate(err, "check email failed")
	}
	if exists {
		return nil, stacktrace.NewError("email already taken")
	}

	hashed, err := password.HashPassword(input.Password)
	if err != nil {
		return nil, stacktrace.Propagate(err, "hash password failed")
	}

	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashed,
		Role:     input.Role,
	}
	model.UserPrepare(user)

	if err := repo.Create(user); err != nil {
		return nil, stacktrace.Propagate(err, "create user failed")
	}

	return user, nil
}

func (d *userDomain) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := d.db.User().FindByID(id)
	if err != nil {
		return nil, stacktrace.Propagate(err, "find user failed")
	}
	if user == nil {
		return nil, stacktrace.NewError("user not found")
	}
	return user, nil
}

func (d *userDomain) GetAll(ctx context.Context, filters model.UserFilter, page, limit int, sort string) ([]model.User, int64, error) {
	return d.db.User().FindAll(filters, page, limit, sort)
}

func (d *userDomain) Update(ctx context.Context, id string, input model.UserInput) (*model.User, error) {
	repo := d.db.User()
	user, err := repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, stacktrace.NewError("user not found")
	}

	if input.Email != "" && input.Email != user.Email {
		exists, _ := repo.ExistsByEmail(input.Email)
		if exists {
			return nil, stacktrace.NewError("email already taken")
		}
		user.Email = input.Email
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	if input.Password != "" {
		hashed, err := password.HashPassword(input.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashed
	}

	user.UpdatedAt = time.Now()
	if err := repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (d *userDomain) Delete(ctx context.Context, id string) error {
	repo := d.db.User()
	if _, err := repo.FindByID(id); err != nil {
		return stacktrace.NewError("user not found")
	}
	return repo.Delete(id)
}