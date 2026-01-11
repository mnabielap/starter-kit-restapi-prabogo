package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Email           string    `json:"email" db:"email"`
	Password        string    `json:"-" db:"password"` // "-" prevents returning in JSON
	Role            string    `json:"role" db:"role"`
	IsEmailVerified bool      `json:"is_email_verified" db:"is_email_verified"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type UserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserFilter struct {
	IDs    []string
	Emails []string
	Role   string
	Search string // For fuzzy search on name/email
}

func UserPrepare(u *User) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
	if u.Role == "" {
		u.Role = "user"
	}
}