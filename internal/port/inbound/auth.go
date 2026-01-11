package inbound_port

type AuthHttpPort interface {
	Register(a any) error
	Login(a any) error
	RefreshToken(a any) error
	Logout(a any) error
	ForgotPassword(a any) error
	ResetPassword(a any) error
}