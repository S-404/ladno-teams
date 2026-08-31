package service

import (
	"fmt"
	"ladno-teams/internal/entity"
	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/repository"

	"github.com/google/uuid"
)

type IProfileService interface {
	Create(userGuid uuid.UUID, name string, about *string) (*entity.Profile, *exception.ApiError)
	GetByUserGuid(userGuid uuid.UUID) (*entity.Profile, *exception.ApiError)
	Update(userGuid uuid.UUID, req dto.ProfileUpdateRequestDto) (*entity.Profile, *exception.ApiError)
	DeleteByUserGuid(userGuid uuid.UUID) *exception.ApiError
}

type ProfileService struct {
	BaseService
	profileRepository repository.IProfileRepository
}

func NewProfileService(baseService *BaseService, profileRepository repository.IProfileRepository) *ProfileService {
	return &ProfileService{
		BaseService:       *baseService,
		profileRepository: profileRepository,
	}
}

func (s *ProfileService) Create(userGuid uuid.UUID, name string, about *string) (*entity.Profile, *exception.ApiError) {
	profile, err := s.profileRepository.Create(entity.Profile{
		UserGuid: userGuid,
		Name:     name,
		About:    about,
	})
	if err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed profile creating: %s", err.Error()))
	}
	return profile, nil
}

func (s *ProfileService) GetByUserGuid(userGuid uuid.UUID) (*entity.Profile, *exception.ApiError) {
	profile, err := s.profileRepository.FindByUserGuid(userGuid)
	if err != nil {
		return nil, exception.EntityNotFoundError("profile", fmt.Sprintf("user_guid:%s", userGuid))
	}
	return profile, nil
}

func (s *ProfileService) Update(userGuid uuid.UUID, req dto.ProfileUpdateRequestDto) (*entity.Profile, *exception.ApiError) {
	profile, apiErr := s.GetByUserGuid(userGuid)
	if apiErr != nil {
		profile = &entity.Profile{UserGuid: userGuid}
	}

	profile.Name = req.Name
	profile.About = req.About

	if apiErr == nil {
		updated, err := s.profileRepository.Update(*profile)
		if err != nil {
			return nil, exception.InternalError("failed profile update")
		}
		return updated, nil
	}

	created, err := s.profileRepository.Create(*profile)
	if err != nil {
		return nil, exception.InternalError("failed profile creating on update")
	}
	return created, nil
}

func (s *ProfileService) DeleteByUserGuid(userGuid uuid.UUID) *exception.ApiError {
	if err := s.profileRepository.DeleteByUserGuid(userGuid); err != nil {
		return exception.InternalError("failed profile delete")
	}
	return nil
}
