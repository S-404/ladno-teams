package entity

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	UserGuid  uuid.UUID `db:"user_guid" json:"user_guid"`
	Name      string    `db:"name" json:"name"`
	About     *string   `db:"about" json:"about,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
