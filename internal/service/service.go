package service

import (
	"ladno-teams/internal/config"
	"ladno-teams/internal/repository"
)

type Service struct {
	Auth          IAuthService
	User          IUserService
	Session       ISessionService
	Profile       IProfileService
	Team          ITeamService
	Teammate      ITeammateService
	Invite        IInviteService
	Workspace     IWorkspaceService
	WorkspaceRole IWorkspaceRoleService
}

func NewService(cfg config.Config, repo *repository.Repository) *Service {
	baseService := NewBaseService(cfg)

	userService := NewUserService(baseService, repo.User)
	profileService := NewProfileService(baseService, repo.Profile)
	sessionService := NewSessionService(baseService, repo.Session)
	authService := NewAuthService(baseService, sessionService)
	teamService := NewTeamService(baseService, repo.Team, repo.Teammate)
	teammateService := NewTeammateService(baseService, repo.Teammate)
	inviteService := NewInviteService(baseService, repo.Invite, repo.Team, repo.Teammate, userService, profileService)
	workspaceService := NewWorkspaceService(baseService, repo.Workspace, repo.WorkspaceRole, repo.Teammate)
	workspaceRoleService := NewWorkspaceRoleService(baseService, repo.WorkspaceRole, repo.Workspace, repo.Teammate)

	return &Service{
		Auth:          authService,
		User:          userService,
		Session:       sessionService,
		Profile:       profileService,
		Team:          teamService,
		Teammate:      teammateService,
		Invite:        inviteService,
		Workspace:     workspaceService,
		WorkspaceRole: workspaceRoleService,
	}
}
