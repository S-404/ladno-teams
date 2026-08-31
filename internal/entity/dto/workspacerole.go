package dto

import (
	"ladno-teams/internal/entity"

	"github.com/google/uuid"
)

type WorkspaceRoleCreateRequestDto struct {
	TeammateGuid uuid.UUID   `json:"teammate_guid" binding:"required"`
	Role         entity.Role `json:"role" binding:"required,valid_role"`
}

type WorkspaceRoleUpdateRequestDto struct {
	Role entity.Role `json:"role" binding:"required,valid_role"`
}

type WorkspaceRoleResponseDto struct {
	WorkspaceGuid uuid.UUID   `json:"workspace_guid"`
	TeammateGuid  uuid.UUID   `json:"teammate_guid"`
	Role          entity.Role `json:"role"`
	CreatedAt     string      `json:"created_at"`
}
