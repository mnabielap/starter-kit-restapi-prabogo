package inbound_port

type UserHttpPort interface {
	Create(a any) error
	GetList(a any) error
	GetOne(a any) error
	Update(a any) error
	Delete(a any) error
}