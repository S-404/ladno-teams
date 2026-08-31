package service

import (
	"fmt"
	"ladno-teams/internal/entity"
	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/repository"

	"github.com/google/uuid"
)

type IWorkspaceService interface {
	Create(userGuid, teamGuid uuid.UUID, req dto.WorkspaceCreateRequestDto) (*entity.Workspace, *exception.ApiError)
	GetByGuid(userGuid, workspaceGuid uuid.UUID) (*entity.Workspace, *exception.ApiError)
	Update(userGuid, workspaceGuid uuid.UUID, req dto.WorkspaceUpdateRequestDto) (*entity.Workspace, *exception.ApiError)
	Delete(userGuid, workspaceGuid uuid.UUID) *exception.ApiError
	AdminDelete(workspaceGuid uuid.UUID) *exception.ApiError
	AdminCreate(teamGuid uuid.UUID, name string) (*entity.Workspace, *exception.ApiError)
	ListByTeam(userGuid, teamGuid uuid.UUID, limit, offset int) ([]entity.Workspace, *exception.ApiError)
	ListAll(limit, offset int) ([]entity.Workspace, *exception.ApiError)
}

type WorkspaceService struct {
	BaseService
	workspaceRepository     repository.IWorkspaceRepository
	workspaceRoleRepository repository.IWorkspaceRoleRepository
	teammateRepository      repository.ITeammateRepository
}

func NewWorkspaceService(
	baseService *BaseService,
	workspaceRepository repository.IWorkspaceRepository,
	workspaceRoleRepository repository.IWorkspaceRoleRepository,
	teammateRepository repository.ITeammateRepository,
) *WorkspaceService {
	return &WorkspaceService{
		BaseService:             *baseService,
		workspaceRepository:     workspaceRepository,
		workspaceRoleRepository: workspaceRoleRepository,
		teammateRepository:      teammateRepository,
	}
}

func (s *WorkspaceService) Create(userGuid, teamGuid uuid.UUID, req dto.WorkspaceCreateRequestDto) (*entity.Workspace, *exception.ApiError) {
	if err := s.requireLeader(userGuid, teamGuid); err != nil {
		return nil, err
	}

	version := req.Version
	if version == "" {
		version = "1.0.0"
	}

	workspace, err := s.workspaceRepository.Create(entity.Workspace{
		TeamGuid: teamGuid,
		Name:     req.Name,
		Version:  version,
		Data:     entity.JSONMapFromDTO(req.Data),
		Config:   entity.JSONMapFromDTO(req.Config),
		Envs:     entity.JSONMapFromDTO(req.Envs),
	})
	if err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed workspace creating: %s", err.Error()))
	}
	return workspace, nil
}

func (s *WorkspaceService) GetByGuid(userGuid, workspaceGuid uuid.UUID) (*entity.Workspace, *exception.ApiError) {
	workspace, err := s.workspaceRepository.FindByGuid(workspaceGuid)
	if err != nil {
		return nil, exception.EntityNotFoundError("workspace", fmt.Sprintf("guid:%s", workspaceGuid))
	}

	if err := s.requireReadAccess(userGuid, workspace); err != nil {
		return nil, err
	}

	return workspace, nil
}

func (s *WorkspaceService) Update(userGuid, workspaceGuid uuid.UUID, req dto.WorkspaceUpdateRequestDto) (*entity.Workspace, *exception.ApiError) {
	workspace, apiErr := s.GetByGuid(userGuid, workspaceGuid)
	if apiErr != nil {
		return nil, apiErr
	}

	if err := s.requireWriteAccess(userGuid, workspace); err != nil {
		return nil, err
	}

	if req.Name != "" {
		workspace.Name = req.Name
	}
	if req.Version != "" {
		workspace.Version = req.Version
	}
	if req.Data != nil {
		workspace.Data = entity.JSONMapFromDTO(req.Data)
	}
	if req.Config != nil {
		workspace.Config = entity.JSONMapFromDTO(req.Config)
	}
	if req.Envs != nil {
		workspace.Envs = entity.JSONMapFromDTO(req.Envs)
	}

	updated, err := s.workspaceRepository.Update(*workspace)
	if err != nil {
		return nil, exception.InternalError("failed workspace update")
	}
	return updated, nil
}

func (s *WorkspaceService) Delete(userGuid, workspaceGuid uuid.UUID) *exception.ApiError {
	workspace, apiErr := s.GetByGuid(userGuid, workspaceGuid)
	if apiErr != nil {
		return apiErr
	}

	if err := s.requireManageAccess(userGuid, workspace); err != nil {
		return err
	}

	if err := s.workspaceRepository.Delete(workspaceGuid); err != nil {
		return exception.InternalError("failed workspace delete")
	}
	return nil
}

func (s *WorkspaceService) AdminDelete(workspaceGuid uuid.UUID) *exception.ApiError {
	if _, err := s.workspaceRepository.FindByGuid(workspaceGuid); err != nil {
		return exception.EntityNotFoundError("workspace", fmt.Sprintf("guid:%s", workspaceGuid))
	}
	if err := s.workspaceRepository.Delete(workspaceGuid); err != nil {
		return exception.InternalError("failed workspace delete")
	}
	return nil
}

func (s *WorkspaceService) AdminCreate(teamGuid uuid.UUID, name string) (*entity.Workspace, *exception.ApiError) {
	if name == "" {
		return nil, exception.BadRequest("name is required")
	}

	emptyJSON := entity.JSONMap{}
	workspace, err := s.workspaceRepository.Create(entity.Workspace{
		TeamGuid: teamGuid,
		Name:     name,
		Version:  "v0",
		Data:     emptyJSON,
		Config:   emptyJSON,
		Envs:     emptyJSON,
	})
	if err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed workspace creating: %s", err.Error()))
	}
	return workspace, nil
}

func (s *WorkspaceService) ListByTeam(userGuid, teamGuid uuid.UUID, limit, offset int) ([]entity.Workspace, *exception.ApiError) {
	_, err := s.teammateRepository.FindByUserAndTeam(userGuid, teamGuid)
	if err != nil {
		return nil, exception.Forbidden("team membership required")
	}

	workspaces, err := s.workspaceRepository.ListByTeam(teamGuid, limit, offset)
	if err != nil {
		return nil, exception.InternalError("failed list workspaces")
	}
	return workspaces, nil
}

func (s *WorkspaceService) ListAll(limit, offset int) ([]entity.Workspace, *exception.ApiError) {
	workspaces, err := s.workspaceRepository.ListAll(limit, offset)
	if err != nil {
		return nil, exception.InternalError("failed list workspaces")
	}
	return workspaces, nil
}

func (s *WorkspaceService) requireLeader(userGuid, teamGuid uuid.UUID) *exception.ApiError {
	isLeader, err := s.teammateRepository.IsLeader(userGuid, teamGuid)
	if err != nil || !isLeader {
		return exception.Forbidden("leader access required")
	}
	return nil
}

func (s *WorkspaceService) requireReadAccess(userGuid uuid.UUID, workspace *entity.Workspace) *exception.ApiError {
	isLeader, err := s.teammateRepository.IsLeader(userGuid, workspace.TeamGuid)
	if err != nil {
		return exception.InternalError("failed leader check")
	}
	if isLeader {
		return nil
	}

	_, err = s.workspaceRoleRepository.FindByWorkspaceAndTeammate(workspace.Guid, userGuid)
	if err != nil {
		return exception.Forbidden("workspace access required")
	}
	return nil
}

func (s *WorkspaceService) requireWriteAccess(userGuid uuid.UUID, workspace *entity.Workspace) *exception.ApiError {
	role, apiErr := s.getEffectiveRole(userGuid, workspace)
	if apiErr != nil {
		return apiErr
	}
	if role == entity.RoleGuest {
		return exception.Forbidden("write access required")
	}
	return nil
}

func (s *WorkspaceService) requireManageAccess(userGuid uuid.UUID, workspace *entity.Workspace) *exception.ApiError {
	role, apiErr := s.getEffectiveRole(userGuid, workspace)
	if apiErr != nil {
		return apiErr
	}
	if role != entity.RoleMaintainer && !s.isLeader(userGuid, workspace.TeamGuid) {
		return exception.Forbidden("manage access required")
	}
	return nil
}

func (s *WorkspaceService) getEffectiveRole(userGuid uuid.UUID, workspace *entity.Workspace) (entity.Role, *exception.ApiError) {
	if s.isLeader(userGuid, workspace.TeamGuid) {
		return entity.RoleMaintainer, nil
	}
	wr, err := s.workspaceRoleRepository.FindByWorkspaceAndTeammate(workspace.Guid, userGuid)
	if err != nil {
		return entity.RoleGuest, exception.Forbidden("workspace access required")
	}
	return wr.Role, nil
}

func (s *WorkspaceService) isLeader(userGuid, teamGuid uuid.UUID) bool {
	isLeader, err := s.teammateRepository.IsLeader(userGuid, teamGuid)
	if err != nil {
		return false
	}
	return isLeader
}
