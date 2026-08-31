package service

import (
	"fmt"
	"ladno-teams/internal/entity"
	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/repository"

	"github.com/google/uuid"
)

type ITeamService interface {
	Create(userGuid uuid.UUID, req dto.TeamCreateRequestDto) (*entity.Team, *exception.ApiError)
	GetByGuid(guid uuid.UUID) (*entity.Team, *exception.ApiError)
	Update(userGuid, teamGuid uuid.UUID, req dto.TeamUpdateRequestDto) (*entity.Team, *exception.ApiError)
	Delete(userGuid, teamGuid uuid.UUID) *exception.ApiError
	AdminDelete(guid uuid.UUID) *exception.ApiError
	ListByUser(userGuid uuid.UUID, limit, offset int) ([]entity.Team, *exception.ApiError)
	ListAll(limit, offset int) ([]entity.Team, *exception.ApiError)
}

type TeamService struct {
	BaseService
	teamRepository     repository.ITeamRepository
	teammateRepository repository.ITeammateRepository
}

func NewTeamService(baseService *BaseService, teamRepository repository.ITeamRepository, teammateRepository repository.ITeammateRepository) *TeamService {
	return &TeamService{
		BaseService:        *baseService,
		teamRepository:     teamRepository,
		teammateRepository: teammateRepository,
	}
}

func (s *TeamService) Create(userGuid uuid.UUID, req dto.TeamCreateRequestDto) (*entity.Team, *exception.ApiError) {
	team, err := s.teamRepository.Create(entity.Team{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed team creating: %s", err.Error()))
	}

	if _, err := s.teammateRepository.Create(entity.Teammate{
		UserGuid: userGuid,
		TeamGuid: team.Guid,
		IsLeader: true,
	}); err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed teammate creating: %s", err.Error()))
	}

	return team, nil
}

func (s *TeamService) GetByGuid(guid uuid.UUID) (*entity.Team, *exception.ApiError) {
	team, err := s.teamRepository.FindByGuid(guid)
	if err != nil {
		return nil, exception.EntityNotFoundError("team", fmt.Sprintf("guid:%s", guid))
	}
	return team, nil
}

func (s *TeamService) Update(userGuid, teamGuid uuid.UUID, req dto.TeamUpdateRequestDto) (*entity.Team, *exception.ApiError) {
	if err := s.requireLeader(userGuid, teamGuid); err != nil {
		return nil, err
	}

	team, apiErr := s.GetByGuid(teamGuid)
	if apiErr != nil {
		return nil, apiErr
	}

	if req.Name != "" {
		team.Name = req.Name
	}
	if req.Description != nil {
		team.Description = req.Description
	}

	updated, err := s.teamRepository.Update(*team)
	if err != nil {
		return nil, exception.InternalError("failed team update")
	}
	return updated, nil
}

func (s *TeamService) Delete(userGuid, teamGuid uuid.UUID) *exception.ApiError {
	if err := s.requireLeader(userGuid, teamGuid); err != nil {
		return err
	}

	if _, apiErr := s.GetByGuid(teamGuid); apiErr != nil {
		return apiErr
	}

	if err := s.teamRepository.Delete(teamGuid); err != nil {
		return exception.InternalError("failed team delete")
	}
	return nil
}

func (s *TeamService) AdminDelete(guid uuid.UUID) *exception.ApiError {
	if _, apiErr := s.GetByGuid(guid); apiErr != nil {
		return apiErr
	}
	if err := s.teamRepository.Delete(guid); err != nil {
		return exception.InternalError("failed team delete")
	}
	return nil
}

func (s *TeamService) ListByUser(userGuid uuid.UUID, limit, offset int) ([]entity.Team, *exception.ApiError) {
	teams, err := s.teamRepository.ListByUser(userGuid, limit, offset)
	if err != nil {
		return nil, exception.InternalError("failed list teams")
	}
	return teams, nil
}

func (s *TeamService) ListAll(limit, offset int) ([]entity.Team, *exception.ApiError) {
	teams, err := s.teamRepository.ListAll(limit, offset)
	if err != nil {
		return nil, exception.InternalError("failed list teams")
	}
	return teams, nil
}

func (s *TeamService) requireLeader(userGuid, teamGuid uuid.UUID) *exception.ApiError {
	isLeader, err := s.teammateRepository.IsLeader(userGuid, teamGuid)
	if err != nil || !isLeader {
		return exception.Forbidden("leader access required")
	}
	return nil
}
