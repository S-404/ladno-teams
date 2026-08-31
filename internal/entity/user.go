package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Guid      uuid.UUID  `db:"guid" json:"guid"`
	Login     string     `db:"login" json:"login"`
	Password  string     `db:"password" json:"-"`
	IsAdmin   bool       `db:"is_admin" json:"is_admin"`
	IsBlocked bool       `db:"is_blocked" json:"is_blocked"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
