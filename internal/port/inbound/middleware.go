package inbound_port

type MiddlewareHttpPort interface {
	Auth(a any) error
	RequireAdmin(a any) error
	RequireAdminOrSelf(a any) error

	InternalAuth(a any) error
	ClientAuth(a any) error
}