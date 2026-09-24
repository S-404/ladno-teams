package dto

import "github.com/google/uuid"

type WorkspaceCreateRequestDto struct {
	Name    string                 `json:"name" binding:"required,max=255"`
	Version string                 `json:"version,omitempty" binding:"omitempty,max=50"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Config  map[string]interface{} `json:"config,omitempty"`
	Envs    map[string]interface{} `json:"envs,omitempty"`
}

type WorkspaceUpdateRequestDto struct {
	Name    string                 `json:"name,omitempty" binding:"omitempty,max=255"`
	Version string                 `json:"version,omitempty" binding:"omitempty,max=50"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Config  map[string]interface{} `json:"config,omitempty"`
	Envs    map[string]interface{} `json:"envs,omitempty"`
}

type WorkspaceListResponseDto struct {
	Guid      uuid.UUID `json:"guid"`
	TeamGuid  uuid.UUID `json:"team_guid"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	GitURL    string    `json:"git_url,omitempty"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
