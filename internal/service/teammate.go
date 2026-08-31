package service

import (
	"fmt"
	"ladno-teams/internal/entity"
	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/repository"

	"github.com/google/uuid"
)

type ITeammateService interface {
	Create(userGuid, teamGuid uuid.UUID, isLeader bool) (*entity.Teammate, *exception.ApiError)
	ListByTeam(teamGuid uuid.UUID) ([]entity.Teammate, *exception.ApiError)
	Update(actorGuid, teamGuid, targetUserGuid uuid.UUID, req dto.TeammateUpdateRequestDto) (*entity.Teammate, *exception.ApiError)
	Delete(actorGuid, teamGuid, targetUserGuid uuid.UUID) *exception.ApiError
	AdminDelete(userGuid, teamGuid uuid.UUID) *exception.ApiError
	IsLeader(userGuid, teamGuid uuid.UUID) (bool, *exception.ApiError)
	IsMember(userGuid, teamGuid uuid.UUID) (bool, *exception.ApiError)
}

type TeammateService struct {
	BaseService
	teammateRepository repository.ITeammateRepository
}

func NewTeammateService(baseService *BaseService, teammateRepository repository.ITeammateRepository) *TeammateService {
	return &TeammateService{
		BaseService:        *baseService,
		teammateRepository: teammateRepository,
	}
}

func (s *TeammateService) Create(userGuid, teamGuid uuid.UUID, isLeader bool) (*entity.Teammate, *exception.ApiError) {
	teammate, err := s.teammateRepository.Create(entity.Teammate{
		UserGuid: userGuid,
		TeamGuid: teamGuid,
		IsLeader: isLeader,
	})
	if err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed teammate creating: %s", err.Error()))
	}
	return teammate, nil
}

func (s *TeammateService) ListByTeam(teamGuid uuid.UUID) ([]entity.Teammate, *exception.ApiError) {
	teammates, err := s.teammateRepository.FindByTeam(teamGuid)
	if err != nil {
		return nil, exception.InternalError("failed list teammates")
	}
	return teammates, nil
}

func (s *TeammateService) Update(actorGuid, teamGuid, targetUserGuid uuid.UUID, req dto.TeammateUpdateRequestDto) (*entity.Teammate, *exception.ApiError) {
	if err := s.requireLeader(actorGuid, teamGuid); err != nil {
		return nil, err
	}

	if _, err := s.teammateRepository.FindByUserAndTeam(targetUserGuid, teamGuid); err != nil {
		return nil, exception.EntityNotFoundError("teammate", fmt.Sprintf("user:%s team:%s", targetUserGuid, teamGuid))
	}

	updated, err := s.teammateRepository.Update(entity.Teammate{
		UserGuid: targetUserGuid,
		TeamGuid: teamGuid,
		IsLeader: req.IsLeader,
	})
	if err != nil {
		return nil, exception.InternalError("failed teammate update")
	}
	return updated, nil
}

func (s *TeammateService) Delete(actorGuid, teamGuid, targetUserGuid uuid.UUID) *exception.ApiError {
	if err := s.requireLeader(actorGuid, teamGuid); err != nil {
		return err
	}

	if actorGuid == targetUserGuid {
		isLeader, err := s.teammateRepository.IsLeader(actorGuid, teamGuid)
		if err != nil {
			return exception.InternalError("failed leader check")
		}
		hasLeader, err := s.teammateRepository.HasLeader(teamGuid)
		if err != nil {
			return exception.InternalError("failed leader check")
		}
		if isLeader && hasLeader {
			count, err := s.teammateRepository.FindByTeam(teamGuid)
			if err != nil {
				return exception.InternalError("failed team count")
			}
			leaders := 0
			for _, t := range count {
				if t.IsLeader {
					leaders++
				}
			}
			if leaders <= 1 {
				return exception.BadRequest("cannot remove the last leader")
			}
		}
	}

	if err := s.teammateRepository.Delete(targetUserGuid, teamGuid); err != nil {
		return exception.InternalError("failed teammate delete")
	}
	return nil
}

func (s *TeammateService) AdminDelete(userGuid, teamGuid uuid.UUID) *exception.ApiError {
	if _, err := s.teammateRepository.FindByUserAndTeam(userGuid, teamGuid); err != nil {
		return exception.EntityNotFoundError("teammate", fmt.Sprintf("user:%s team:%s", userGuid, teamGuid))
	}
	if err := s.teammateRepository.Delete(userGuid, teamGuid); err != nil {
		return exception.InternalError("failed teammate delete")
	}
	return nil
}

func (s *TeammateService) IsLeader(userGuid, teamGuid uuid.UUID) (bool, *exception.ApiError) {
	isLeader, err := s.teammateRepository.IsLeader(userGuid, teamGuid)
	if err != nil {
		return false, exception.InternalError("failed leader check")
	}
	return isLeader, nil
}

func (s *TeammateService) IsMember(userGuid, teamGuid uuid.UUID) (bool, *exception.ApiError) {
	isMember, err := s.teammateRepository.IsMember(userGuid, teamGuid)
	if err != nil {
		return false, exception.InternalError("failed member check")
	}
	return isMember, nil
}

func (s *TeammateService) requireLeader(userGuid, teamGuid uuid.UUID) *exception.ApiError {
	isLeader, err := s.teammateRepository.IsLeader(userGuid, teamGuid)
	if err != nil || !isLeader {
		return exception.Forbidden("leader access required")
	}
	return nil
}
