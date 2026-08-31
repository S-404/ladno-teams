package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleMaintainer Role = "MAINTAINER"
	RoleDeveloper  Role = "DEVELOPER"
	RoleGuest      Role = "GUEST"
)

type WorkspaceRole struct {
	WorkspaceGuid uuid.UUID `db:"workspace_guid" json:"workspace_guid"`
	TeammateGuid  uuid.UUID `db:"teammate_guid" json:"teammate_guid"`
	Role          Role      `db:"role" json:"role"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
