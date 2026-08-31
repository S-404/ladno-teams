package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	User          IUserRepository
	Session       ISessionRepository
	Profile       IProfileRepository
	Team          ITeamRepository
	Teammate      ITeammateRepository
	Invite        IInviteRepository
	Workspace     IWorkspaceRepository
	WorkspaceRole IWorkspaceRoleRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	baseRepo := NewBaseRepository(db)

	return &Repository{
		User:          NewUserRepository(baseRepo),
		Session:       NewSessionRepository(baseRepo),
		Profile:       NewProfileRepository(baseRepo),
		Team:          NewTeamRepository(baseRepo),
		Teammate:      NewTeammateRepository(baseRepo),
		Invite:        NewInviteRepository(baseRepo),
		Workspace:     NewWorkspaceRepository(baseRepo),
		WorkspaceRole: NewWorkspaceRoleRepository(baseRepo),
	}
}
