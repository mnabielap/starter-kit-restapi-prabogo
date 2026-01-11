package domain

import (
	"prabogo/internal/domain/auth"
	"prabogo/internal/domain/client"
	"prabogo/internal/domain/user"
	outbound_port "prabogo/internal/port/outbound"
)

type Domain interface {
	Client() client.ClientDomain
	User() user.UserDomain
	Auth() auth.AuthDomain
}

type domain struct {
	databasePort outbound_port.DatabasePort
	messagePort  outbound_port.MessagePort
	cachePort    outbound_port.CachePort
	emailPort    outbound_port.EmailPort
}

func NewDomain(
	databasePort outbound_port.DatabasePort,
	messagePort outbound_port.MessagePort,
	cachePort outbound_port.CachePort,
	emailPort outbound_port.EmailPort,
) Domain {
	return &domain{
		databasePort: databasePort,
		messagePort:  messagePort,
		cachePort:    cachePort,
		emailPort:    emailPort,
	}
}

func (d *domain) Client() client.ClientDomain {
	return client.NewClientDomain(d.databasePort, d.messagePort, d.cachePort)
}

func (d *domain) User() user.UserDomain {
	return user.NewUserDomain(d.databasePort)
}

func (d *domain) Auth() auth.AuthDomain {
	return auth.NewAuthDomain(d.databasePort, d.emailPort)
}