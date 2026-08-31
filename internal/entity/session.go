package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Guid      uuid.UUID `db:"guid" json:"guid"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"token"`
	ExpiredAt time.Time `db:"expired_at" json:"expired_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type SessionWithUser struct {
	Session
	User User `db:"user" json:"user"`
}
