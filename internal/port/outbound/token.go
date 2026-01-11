package outbound_port

import "prabogo/internal/model"

//go:generate mockgen -source=token.go -destination=./../../../tests/mocks/port/mock_token.go
type TokenDatabasePort interface {
	Create(token *model.Token) error
	FindByToken(tokenStr string, tokenType string) (*model.Token, error)
	DeleteByUserIDAndType(userID string, tokenType string) error
	Delete(tokenID int) error
}