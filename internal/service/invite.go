package service

import (
	"fmt"
	"ladno-teams/internal/entity"
	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/repository"
	"time"

	"github.com/google/uuid"
)

type IInviteService interface {
	Create(teamGuid, leaderUserGuid uuid.UUID) (*entity.Invite, *exception.ApiError)
	AdminCreate(teamGuid uuid.UUID) (*entity.Invite, *exception.ApiError)
	GetByGuid(guid uuid.UUID) (*entity.Invite, *exception.ApiError)
	Accept(inviteGuid, userGuid uuid.UUID) *exception.ApiError
	RegisterByInvite(req dto.InviteRegisterRequestDto) (*entity.User, *exception.ApiError)
	Delete(guid uuid.UUID) *exception.ApiError
}

type InviteService struct {
	BaseService
	inviteRepository   repository.IInviteRepository
	teamRepository     repository.ITeamRepository
	teammateRepository repository.ITeammateRepository
	userService        IUserService
	profileService     IProfileService
}

func NewInviteService(
	baseService *BaseService,
	inviteRepository repository.IInviteRepository,
	teamRepository repository.ITeamRepository,
	teammateRepository repository.ITeammateRepository,
	userService IUserService,
	profileService IProfileService,
) *InviteService {
	return &InviteService{
		BaseService:        *baseService,
		inviteRepository:   inviteRepository,
		teamRepository:     teamRepository,
		teammateRepository: teammateRepository,
		userService:        userService,
		profileService:     profileService,
	}
}

func (s *InviteService) AdminCreate(teamGuid uuid.UUID) (*entity.Invite, *exception.ApiError) {
	if _, err := s.teamRepository.FindByGuid(teamGuid); err != nil {
		return nil, exception.EntityNotFoundError("team", fmt.Sprintf("guid:%s", teamGuid))
	}

	invite, err := s.inviteRepository.Create(entity.Invite{
		TeamGuid:  teamGuid,
		ExpiredAt: time.Now().UTC().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed invite creating: %s", err.Error()))
	}
	return invite, nil
}

func (s *InviteService) Create(teamGuid, leaderUserGuid uuid.UUID) (*entity.Invite, *exception.ApiError) {
	isLeader, err := s.teammateRepository.IsLeader(leaderUserGuid, teamGuid)
	if err != nil || !isLeader {
		return nil, exception.Forbidden("leader access required")
	}

	invite, err := s.inviteRepository.Create(entity.Invite{
		TeamGuid:  teamGuid,
		ExpiredAt: time.Now().UTC().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed invite creating: %s", err.Error()))
	}
	return invite, nil
}

func (s *InviteService) GetByGuid(guid uuid.UUID) (*entity.Invite, *exception.ApiError) {
	invite, err := s.inviteRepository.FindByGuid(guid)
	if err != nil {
		return nil, exception.EntityNotFoundError("invite", fmt.Sprintf("guid:%s", guid))
	}
	return invite, nil
}

func (s *InviteService) Accept(inviteGuid, userGuid uuid.UUID) *exception.ApiError {
	invite, apiErr := s.GetByGuid(inviteGuid)
	if apiErr != nil {
		return apiErr
	}

	isMember, err := s.teammateRepository.IsMember(userGuid, invite.TeamGuid)
	if err != nil {
		return exception.InternalError("failed member check")
	}
	if isMember {
		_ = s.inviteRepository.Delete(inviteGuid)
		return exception.BadRequest("user is already a member of the team")
	}

	if _, err := s.teammateRepository.Create(entity.Teammate{
		UserGuid: userGuid,
		TeamGuid: invite.TeamGuid,
		IsLeader: false,
	}); err != nil {
		return exception.InternalError(fmt.Sprintf("failed teammate creating: %s", err.Error()))
	}

	if err := s.inviteRepository.Delete(inviteGuid); err != nil {
		return exception.InternalError("failed invite delete")
	}
	return nil
}

func (s *InviteService) RegisterByInvite(req dto.InviteRegisterRequestDto) (*entity.User, *exception.ApiError) {
	invite, apiErr := s.GetByGuid(req.Guid)
	if apiErr != nil {
		return nil, apiErr
	}

	user, apiErr := s.userService.Create(dto.AuthRequestDto{
		Login:    req.Login,
		Password: req.Password,
	})
	if apiErr != nil {
		return nil, apiErr
	}

	if _, apiErr := s.profileService.Create(user.Guid, req.Name, nil); apiErr != nil {
		return nil, apiErr
	}

	if _, err := s.teammateRepository.Create(entity.Teammate{
		UserGuid: user.Guid,
		TeamGuid: invite.TeamGuid,
		IsLeader: false,
	}); err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed teammate creating: %s", err.Error()))
	}

	if err := s.inviteRepository.Delete(invite.Guid); err != nil {
		return nil, exception.InternalError("failed invite delete")
	}

	return user, nil
}

func (s *InviteService) Delete(guid uuid.UUID) *exception.ApiError {
	if err := s.inviteRepository.Delete(guid); err != nil {
		return exception.InternalError("failed invite delete")
	}
	return nil
}
