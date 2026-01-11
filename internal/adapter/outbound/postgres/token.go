package postgres_outbound_adapter

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

const tableToken = "tokens"

type tokenAdapter struct {
	db outbound_port.DatabaseExecutor
}

func NewTokenAdapter(db outbound_port.DatabaseExecutor) outbound_port.TokenDatabasePort {
	return &tokenAdapter{db: db}
}

func (a *tokenAdapter) Create(token *model.Token) error {
	// Let Postgres generate the Serial ID, so we exclude ID from insert if it's 0
	ds := goqu.Dialect("postgres").Insert(tableToken).Rows(
		goqu.Record{
			"token":       token.Token,
			"user_id":     token.UserID,
			"type":        token.Type,
			"expires":     token.Expires,
			"blacklisted": token.Blacklisted,
			"created_at":  token.CreatedAt,
		},
	)
	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *tokenAdapter) FindByToken(tokenStr string, tokenType string) (*model.Token, error) {
	ds := goqu.Dialect("postgres").From(tableToken).
		Where(goqu.Ex{
			"token":       tokenStr,
			"type":        tokenType,
			"blacklisted": false,
		})

	query, _, err := ds.ToSQL()
	if err != nil {
		return nil, err
	}

	var t model.Token
	err = a.db.QueryRow(query).Scan(
		&t.ID, &t.Token, &t.UserID, &t.Type, &t.Expires, &t.Blacklisted, &t.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (a *tokenAdapter) DeleteByUserIDAndType(userID string, tokenType string) error {
	ds := goqu.Dialect("postgres").Delete(tableToken).
		Where(goqu.Ex{"user_id": userID, "type": tokenType})

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *tokenAdapter) Delete(tokenID int) error {
	ds := goqu.Dialect("postgres").Delete(tableToken).Where(goqu.Ex{"id": tokenID})
	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}
	_, err = a.db.Exec(query)
	return err
}