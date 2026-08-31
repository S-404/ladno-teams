package entity

import (
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	Guid      uuid.UUID `db:"guid" json:"guid"`
	TeamGuid  uuid.UUID `db:"team_guid" json:"team_guid"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	ExpiredAt time.Time `db:"expired_at" json:"expired_at"`
}
