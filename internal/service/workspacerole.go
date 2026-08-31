package service

import (
	"fmt"
	"ladno-teams/internal/entity"
	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/repository"

	"github.com/google/uuid"
)

type IWorkspaceRoleService interface {
	Create(userGuid, workspaceGuid uuid.UUID, req dto.WorkspaceRoleCreateRequestDto) (*entity.WorkspaceRole, *exception.ApiError)
	ListByWorkspace(userGuid, workspaceGuid uuid.UUID) ([]entity.WorkspaceRole, *exception.ApiError)
	Update(userGuid, workspaceGuid, teammateGuid uuid.UUID, req dto.WorkspaceRoleUpdateRequestDto) (*entity.WorkspaceRole, *exception.ApiError)
	Delete(userGuid, workspaceGuid, teammateGuid uuid.UUID) *exception.ApiError
	AdminDelete(workspaceGuid, teammateGuid uuid.UUID) *exception.ApiError
}

type WorkspaceRoleService struct {
	BaseService
	workspaceRoleRepository repository.IWorkspaceRoleRepository
	workspaceRepository     repository.IWorkspaceRepository
	teammateRepository      repository.ITeammateRepository
}

func NewWorkspaceRoleService(
	baseService *BaseService,
	workspaceRoleRepository repository.IWorkspaceRoleRepository,
	workspaceRepository repository.IWorkspaceRepository,
	teammateRepository repository.ITeammateRepository,
) *WorkspaceRoleService {
	return &WorkspaceRoleService{
		BaseService:             *baseService,
		workspaceRoleRepository: workspaceRoleRepository,
		workspaceRepository:     workspaceRepository,
		teammateRepository:      teammateRepository,
	}
}

func (s *WorkspaceRoleService) Create(userGuid, workspaceGuid uuid.UUID, req dto.WorkspaceRoleCreateRequestDto) (*entity.WorkspaceRole, *exception.ApiError) {
	workspace, err := s.workspaceRepository.FindByGuid(workspaceGuid)
	if err != nil {
		return nil, exception.EntityNotFoundError("workspace", fmt.Sprintf("guid:%s", workspaceGuid))
	}

	if err := s.requireRoleManager(userGuid, workspace); err != nil {
		return nil, err
	}

	if _, err := s.teammateRepository.FindByUserAndTeam(req.TeammateGuid, workspace.TeamGuid); err != nil {
		return nil, exception.EntityNotFoundError("teammate", fmt.Sprintf("user:%s team:%s", req.TeammateGuid, workspace.TeamGuid))
	}

	role, err := s.workspaceRoleRepository.Create(entity.WorkspaceRole{
		WorkspaceGuid: workspace.Guid,
		TeammateGuid:  req.TeammateGuid,
		Role:          req.Role,
	})
	if err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed workspace role creating: %s", err.Error()))
	}
	return role, nil
}

func (s *WorkspaceRoleService) ListByWorkspace(userGuid, workspaceGuid uuid.UUID) ([]entity.WorkspaceRole, *exception.ApiError) {
	workspace, err := s.workspaceRepository.FindByGuid(workspaceGuid)
	if err != nil {
		return nil, exception.EntityNotFoundError("workspace", fmt.Sprintf("guid:%s", workspaceGuid))
	}

	isMember, err := s.teammateRepository.IsMember(userGuid, workspace.TeamGuid)
	if err != nil || !isMember {
		return nil, exception.Forbidden("team membership required")
	}

	roles, err := s.workspaceRoleRepository.FindByWorkspace(workspaceGuid)
	if err != nil {
		return nil, exception.InternalError("failed list workspace roles")
	}
	return roles, nil
}

func (s *WorkspaceRoleService) Update(userGuid, workspaceGuid, teammateGuid uuid.UUID, req dto.WorkspaceRoleUpdateRequestDto) (*entity.WorkspaceRole, *exception.ApiError) {
	workspace, err := s.workspaceRepository.FindByGuid(workspaceGuid)
	if err != nil {
		return nil, exception.EntityNotFoundError("workspace", fmt.Sprintf("guid:%s", workspaceGuid))
	}

	if err := s.requireRoleManager(userGuid, workspace); err != nil {
		return nil, err
	}

	role, err := s.workspaceRoleRepository.Update(entity.WorkspaceRole{
		WorkspaceGuid: workspaceGuid,
		TeammateGuid:  teammateGuid,
		Role:          req.Role,
	})
	if err != nil {
		return nil, exception.EntityNotFoundError("workspace_role", fmt.Sprintf("workspace:%s teammate:%s", workspaceGuid, teammateGuid))
	}
	return role, nil
}

func (s *WorkspaceRoleService) Delete(userGuid, workspaceGuid, teammateGuid uuid.UUID) *exception.ApiError {
	workspace, err := s.workspaceRepository.FindByGuid(workspaceGuid)
	if err != nil {
		return exception.EntityNotFoundError("workspace", fmt.Sprintf("guid:%s", workspaceGuid))
	}

	if err := s.requireRoleManager(userGuid, workspace); err != nil {
		return err
	}

	if err := s.workspaceRoleRepository.Delete(workspaceGuid, teammateGuid); err != nil {
		return exception.EntityNotFoundError("workspace_role", fmt.Sprintf("workspace:%s teammate:%s", workspaceGuid, teammateGuid))
	}
	return nil
}

func (s *WorkspaceRoleService) AdminDelete(workspaceGuid, teammateGuid uuid.UUID) *exception.ApiError {
	if _, err := s.workspaceRoleRepository.FindByWorkspaceAndTeammate(workspaceGuid, teammateGuid); err != nil {
		return exception.EntityNotFoundError("workspace_role", fmt.Sprintf("workspace:%s teammate:%s", workspaceGuid, teammateGuid))
	}
	if err := s.workspaceRoleRepository.Delete(workspaceGuid, teammateGuid); err != nil {
		return exception.InternalError("failed workspace role delete")
	}
	return nil
}

func (s *WorkspaceRoleService) requireRoleManager(userGuid uuid.UUID, workspace *entity.Workspace) *exception.ApiError {
	isLeader, err := s.teammateRepository.IsLeader(userGuid, workspace.TeamGuid)
	if err != nil {
		return exception.InternalError("failed leader check")
	}
	if isLeader {
		return nil
	}

	wr, err := s.workspaceRoleRepository.FindByWorkspaceAndTeammate(workspace.Guid, userGuid)
	if err != nil || wr.Role != entity.RoleMaintainer {
		return exception.Forbidden("maintainer access required")
	}
	return nil
}
