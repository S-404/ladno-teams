package entity

import (
	"time"

	"github.com/google/uuid"
)

type Teammate struct {
	UserGuid  uuid.UUID `db:"user_guid" json:"user_guid"`
	TeamGuid  uuid.UUID `db:"team_guid" json:"team_guid"`
	IsLeader  bool      `db:"is_leader" json:"is_leader"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
