package entity

import (
	"time"

	"github.com/google/uuid"
)

type Workspace struct {
	Guid      uuid.UUID  `db:"guid" json:"guid"`
	TeamGuid  uuid.UUID  `db:"team_guid" json:"team_guid"`
	Name      string     `db:"name" json:"name"`
	Version   string     `db:"version" json:"version"`
	Data      JSONMap    `db:"data" json:"data"`
	Config    JSONMap    `db:"config" json:"config"`
	Envs      JSONMap    `db:"envs" json:"envs"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
