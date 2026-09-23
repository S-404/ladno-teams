package dto

import "github.com/google/uuid"

type TeammateUpdateRequestDto struct {
	IsLeader bool `json:"is_leader" binding:"required"`
}

type TeammateProfileDto struct {
	Name string `json:"name"`
}

type TeammateUserDto struct {
	Guid    uuid.UUID           `json:"guid"`
	Login   string              `json:"login"`
	Profile *TeammateProfileDto `json:"profile,omitempty"`
}

type TeammateListResponseDto struct {
	UserGuid  uuid.UUID       `json:"user_guid"`
	TeamGuid  uuid.UUID       `json:"team_guid"`
	IsLeader  bool            `json:"is_leader"`
	CreatedAt string          `json:"created_at"`
	User      TeammateUserDto `json:"user"`
}
