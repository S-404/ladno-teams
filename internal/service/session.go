package service

import (
	"fmt"
	"ladno-teams/internal/entity"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/repository"
	"time"

	"github.com/google/uuid"
)

type ISessionService interface {
	Create(userID uuid.UUID, token string, ttl time.Duration) (*entity.Session, *exception.ApiError)
	FindByToken(token string) (*entity.Session, *exception.ApiError)
	DestroyByToken(token string) *exception.ApiError
	DestroyByUserID(userID uuid.UUID) *exception.ApiError
}

type SessionService struct {
	BaseService
	sessionRepository repository.ISessionRepository
}

func NewSessionService(baseService *BaseService, sessionRepository repository.ISessionRepository) *SessionService {
	return &SessionService{
		BaseService:       *baseService,
		sessionRepository: sessionRepository,
	}
}

func (s *SessionService) Create(userID uuid.UUID, token string, ttl time.Duration) (*entity.Session, *exception.ApiError) {
	created, err := s.sessionRepository.Create(entity.Session{
		UserID:    userID,
		Token:     token,
		ExpiredAt: time.Now().UTC().Add(ttl),
	})
	if err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed session creating: %s", err.Error()))
	}
	return created, nil
}

func (s *SessionService) FindByToken(token string) (*entity.Session, *exception.ApiError) {
	found, err := s.sessionRepository.FindByToken(token)
	if err != nil {
		return nil, exception.AuthError("invalid session")
	}
	return found, nil
}

func (s *SessionService) DestroyByToken(token string) *exception.ApiError {
	if err := s.sessionRepository.DeleteByToken(token); err != nil {
		return exception.InternalError("failed session destroy")
	}
	return nil
}

func (s *SessionService) DestroyByUserID(userID uuid.UUID) *exception.ApiError {
	if err := s.sessionRepository.DeleteByUserID(userID); err != nil {
		return exception.InternalError("failed sessions destroy")
	}
	return nil
}
