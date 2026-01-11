package postgres_outbound_adapter

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/palantir/stacktrace"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

const tableUser = "users"

type userAdapter struct {
	db outbound_port.DatabaseExecutor
}

func NewUserAdapter(db outbound_port.DatabaseExecutor) outbound_port.UserDatabasePort {
	return &userAdapter{db: db}
}

func (a *userAdapter) Create(user *model.User) error {
	ds := goqu.Dialect("postgres").Insert(tableUser).Rows(user)
	query, _, err := ds.ToSQL()
	if err != nil {
		return stacktrace.Propagate(err, "failed to build create user query")
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *userAdapter) FindByEmail(email string) (*model.User, error) {
	ds := goqu.Dialect("postgres").From(tableUser).Where(goqu.Ex{"email": email})
	return a.fetchOne(ds)
}

func (a *userAdapter) FindByID(id string) (*model.User, error) {
	ds := goqu.Dialect("postgres").From(tableUser).Where(goqu.Ex{"id": id})
	return a.fetchOne(ds)
}

func (a *userAdapter) ExistsByEmail(email string) (bool, error) {
	ds := goqu.Dialect("postgres").From(tableUser).Select(goqu.L("1")).Where(goqu.Ex{"email": email})
	query, _, err := ds.ToSQL()
	if err != nil {
		return false, err
	}

	var exists int
	err = a.db.QueryRow(query).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *userAdapter) Update(user *model.User) error {
	ds := goqu.Dialect("postgres").Update(tableUser).
		Set(user).
		Where(goqu.Ex{"id": user.ID})

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *userAdapter) Delete(id string) error {
	ds := goqu.Dialect("postgres").Delete(tableUser).Where(goqu.Ex{"id": id})
	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}
	_, err = a.db.Exec(query)
	return err
}

func (a *userAdapter) FindAll(filters model.UserFilter, page, limit int, sort string) ([]model.User, int64, error) {
	ds := goqu.Dialect("postgres").From(tableUser)

	// Filtering Logic
	if filters.Search != "" {
		pattern := "%" + strings.ToLower(filters.Search) + "%"
		ds = ds.Where(goqu.Or(
			goqu.L("LOWER(name) LIKE ?", pattern),
			goqu.L("LOWER(email) LIKE ?", pattern),
		))
	}
	if filters.Role != "" {
		ds = ds.Where(goqu.Ex{"role": filters.Role})
	}

	// Count Total
	countDs := ds.Select(goqu.COUNT("*"))
	countQuery, _, err := countDs.ToSQL()
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := a.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Sorting
	if sort != "" {
		parts := strings.Split(sort, ":")
		col := parts[0]
		dir := "asc"
		if len(parts) > 1 {
			dir = parts[1]
		}
		if dir == "desc" {
			ds = ds.Order(goqu.I(col).Desc())
		} else {
			ds = ds.Order(goqu.I(col).Asc())
		}
	} else {
		ds = ds.Order(goqu.I("created_at").Desc())
	}

	// Pagination
	offset := (page - 1) * limit
	ds = ds.Limit(uint(limit)).Offset(uint(offset))

	// Execute
	query, _, err := ds.ToSQL()
	if err != nil {
		return nil, 0, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsEmailVerified, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}

// Helper to fetch single user
func (a *userAdapter) fetchOne(ds *goqu.SelectDataset) (*model.User, error) {
	query, _, err := ds.ToSQL()
	if err != nil {
		return nil, err
	}

	var u model.User
	err = a.db.QueryRow(query).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsEmailVerified, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Return nil if not found (Domain layer handles 404)
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}