package dto

import "github.com/google/uuid"

type TeammateUpdateRequestDto struct {
	IsLeader bool `json:"is_leader" binding:"required"`
}

type TeammateListResponseDto struct {
	UserGuid  uuid.UUID `json:"user_guid"`
	TeamGuid  uuid.UUID `json:"team_guid"`
	IsLeader  bool      `json:"is_leader"`
	CreatedAt string    `json:"created_at"`
}
