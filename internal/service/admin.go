package service

import (
	"ladno-teams/internal/entity"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/repository"
)

type IAdminService interface {
	ListUsers(limit, offset int) ([]entity.AdminUserView, *exception.ApiError)
	ListInvites(limit, offset int) ([]entity.AdminInviteView, *exception.ApiError)
	ListTeammates(limit, offset int) ([]entity.AdminTeammateView, *exception.ApiError)
	ListWorkspaces(limit, offset int) ([]entity.AdminWorkspaceView, *exception.ApiError)
	ListWorkspaceRoles(limit, offset int) ([]entity.AdminWorkspaceRoleView, *exception.ApiError)
	GetUserViewByGuid(guid string) (*entity.AdminUserView, *exception.ApiError)
}

type AdminService struct {
	BaseService
	adminRepository repository.IAdminRepository
}

func NewAdminService(
	baseService *BaseService,
	adminRepository repository.IAdminRepository,
) *AdminService {
	return &AdminService{
		BaseService:     *baseService,
		adminRepository: adminRepository,
	}
}

func (s *AdminService) ListUsers(limit, offset int) ([]entity.AdminUserView, *exception.ApiError) {
	rows, err := s.adminRepository.ListUsersWithProfiles(limit, offset)
	if err != nil {
		return nil, exception.InternalError("failed list users")
	}
	return rows, nil
}

func (s *AdminService) ListInvites(limit, offset int) ([]entity.AdminInviteView, *exception.ApiError) {
	rows, err := s.adminRepository.ListInvitesWithTeam(limit, offset)
	if err != nil {
		return nil, exception.InternalError("failed list invites")
	}
	return rows, nil
}

func (s *AdminService) ListTeammates(limit, offset int) ([]entity.AdminTeammateView, *exception.ApiError) {
	rows, err := s.adminRepository.ListTeammatesWithDetails(limit, offset)
	if err != nil {
		return nil, exception.InternalError("failed list teammates")
	}
	return rows, nil
}

func (s *AdminService) ListWorkspaces(limit, offset int) ([]entity.AdminWorkspaceView, *exception.ApiError) {
	rows, err := s.adminRepository.ListWorkspacesWithTeam(limit, offset)
	if err != nil {
		return nil, exception.InternalError("failed list workspaces")
	}
	return rows, nil
}

func (s *AdminService) ListWorkspaceRoles(limit, offset int) ([]entity.AdminWorkspaceRoleView, *exception.ApiError) {
	rows, err := s.adminRepository.ListWorkspaceRolesWithDetails(limit, offset)
	if err != nil {
		return nil, exception.InternalError("failed list workspace roles")
	}
	return rows, nil
}

func (s *AdminService) GetUserViewByGuid(guid string) (*entity.AdminUserView, *exception.ApiError) {
	users, err := s.adminRepository.ListUsersWithProfiles(1000, 0)
	if err != nil {
		return nil, exception.InternalError("failed get user view")
	}
	for _, u := range users {
		if u.Guid.String() == guid {
			return &u, nil
		}
	}
	return nil, exception.EntityNotFoundError("user", guid)
}
