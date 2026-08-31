package entity

import (
	"time"

	"github.com/google/uuid"
)

type AdminUserView struct {
	Guid      uuid.UUID `db:"guid"`
	Login     string    `db:"login"`
	IsAdmin   bool      `db:"is_admin"`
	IsBlocked bool      `db:"is_blocked"`
	Name      *string   `db:"name"`
	About     *string   `db:"about"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type AdminInviteView struct {
	Guid      uuid.UUID `db:"guid"`
	TeamGuid  uuid.UUID `db:"team_guid"`
	TeamName  string    `db:"team_name"`
	ExpiredAt time.Time `db:"expired_at"`
}

type AdminTeammateView struct {
	UserGuid uuid.UUID `db:"user_guid"`
	TeamGuid uuid.UUID `db:"team_guid"`
	IsLeader bool      `db:"is_leader"`
	TeamName string    `db:"team_name"`
	UserName string    `db:"user_name"`
}

type AdminWorkspaceView struct {
	Guid     uuid.UUID `db:"guid"`
	Name     string    `db:"name"`
	Version  string    `db:"version"`
	TeamName string    `db:"team_name"`
}

type AdminWorkspaceRoleView struct {
	WorkspaceGuid uuid.UUID `db:"workspace_guid"`
	TeammateGuid  uuid.UUID `db:"teammate_guid"`
	Role          Role      `db:"role"`
	TeamName      string    `db:"team_name"`
	WorkspaceName string    `db:"workspace_name"`
	UserName      string    `db:"user_name"`
}
