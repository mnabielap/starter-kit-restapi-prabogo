package model

import (
	"time"
)

const (
	TokenTypeAccess        = "access"
	TokenTypeRefresh       = "refresh"
	TokenTypeResetPassword = "resetPassword"
	TokenTypeVerifyEmail   = "verifyEmail"
)

type Token struct {
	ID          int       `json:"id" db:"id"`
	Token       string    `json:"token" db:"token"`
	UserID      string    `json:"user_id" db:"user_id"`
	Type        string    `json:"type" db:"type"`
	Expires     time.Time `json:"expires" db:"expires"`
	Blacklisted bool      `json:"blacklisted" db:"blacklisted"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type TokenInput struct {
	Token   string
	UserID  string
	Type    string
	Expires time.Time
}