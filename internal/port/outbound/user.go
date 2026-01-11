package outbound_port

import "prabogo/internal/model"

//go:generate mockgen -source=user.go -destination=./../../../tests/mocks/port/mock_user.go
type UserDatabasePort interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	// FindAll returns users and total count
	FindAll(filters model.UserFilter, page, limit int, sort string) ([]model.User, int64, error)
	ExistsByEmail(email string) (bool, error)
	Update(user *model.User) error
	Delete(id string) error
}