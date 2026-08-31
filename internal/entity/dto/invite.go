package dto

import (
	"time"

	"github.com/google/uuid"
)

type InviteCreateResponseDto struct {
	Guid      uuid.UUID `json:"guid"`
	TeamGuid  uuid.UUID `json:"team_guid"`
	ExpiredAt time.Time `json:"expired_at"`
}

type InviteCheckResponseDto struct {
	Guid      uuid.UUID `json:"guid"`
	TeamGuid  uuid.UUID `json:"team_guid"`
	ExpiredAt time.Time `json:"expired_at"`
	IsExpired bool      `json:"is_expired"`
}
