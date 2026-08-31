package service

import (
	"fmt"
	"ladno-teams/internal/config"
	"ladno-teams/internal/entity"
	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Create(req dto.AuthRequestDto) (*entity.User, *exception.ApiError)
	GetByGuid(guid uuid.UUID) (*entity.User, *exception.ApiError)
	GetByLogin(login string) (*entity.User, *exception.ApiError)
	GetByCredentials(login, password string) (*entity.User, *exception.ApiError)
	Update(guid uuid.UUID, req dto.UserUpdateRequestDto) (*entity.User, *exception.ApiError)
	AdminUpdate(guid uuid.UUID, req dto.UserAdminUpdateRequestDto) (*entity.User, *exception.ApiError)
	Delete(guid uuid.UUID) *exception.ApiError
	List(limit, offset int) ([]entity.User, *exception.ApiError)
	EnsureAdmin(cfg config.Config) *exception.ApiError
}

type UserService struct {
	BaseService
	userRepository repository.IUserRepository
}

func NewUserService(baseService *BaseService, userRepository repository.IUserRepository) *UserService {
	return &UserService{
		BaseService:    *baseService,
		userRepository: userRepository,
	}
}

func (s *UserService) Create(req dto.AuthRequestDto) (*entity.User, *exception.ApiError) {
	candidate, err := s.userRepository.FindByLogin(req.Login)
	if err == nil && candidate.Login == req.Login {
		return nil, exception.EntityAlreadyExistsError("user", fmt.Sprintf("login:%s", req.Login))
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, exception.InternalError("failed password hashing: create user")
	}

	createdUser, err := s.userRepository.Create(entity.User{
		Login:    req.Login,
		Password: hashedPassword,
	})
	if err != nil {
		return nil, exception.InternalError(fmt.Sprintf("failed user creating with login '%s'", req.Login))
	}

	return createdUser, nil
}

func (s *UserService) GetByLogin(login string) (*entity.User, *exception.ApiError) {
	user, err := s.userRepository.FindByLogin(login)
	if err != nil {
		return nil, exception.EntityNotFoundError("user", fmt.Sprintf("login:%s", login))
	}
	return user, nil
}

func (s *UserService) GetByCredentials(login, password string) (*entity.User, *exception.ApiError) {
	user, apiErr := s.GetByLogin(login)
	if apiErr != nil {
		return nil, apiErr
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, exception.BadRequest("bad credentials")
	}

	return user, nil
}

func (s *UserService) Update(guid uuid.UUID, req dto.UserUpdateRequestDto) (*entity.User, *exception.ApiError) {
	user, apiErr := s.GetByGuid(guid)
	if apiErr != nil {
		return nil, apiErr
	}

	if req.Login != "" && req.Login != user.Login {
		if _, err := s.userRepository.FindByLogin(req.Login); err == nil {
			return nil, exception.EntityAlreadyExistsError("user", fmt.Sprintf("login:%s", req.Login))
		}
		user.Login = req.Login
	}

	updated, err := s.userRepository.Update(*user)
	if err != nil {
		return nil, exception.InternalError("failed user update")
	}
	return updated, nil
}

func (s *UserService) AdminUpdate(guid uuid.UUID, req dto.UserAdminUpdateRequestDto) (*entity.User, *exception.ApiError) {
	user, apiErr := s.GetByGuid(guid)
	if apiErr != nil {
		return nil, apiErr
	}

	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}
	if req.IsBlocked != nil {
		user.IsBlocked = *req.IsBlocked
	}

	updated, err := s.userRepository.Update(*user)
	if err != nil {
		return nil, exception.InternalError("failed user admin update")
	}
	return updated, nil
}

func (s *UserService) Delete(guid uuid.UUID) *exception.ApiError {
	if _, err := s.userRepository.FindByGuid(guid); err != nil {
		return exception.EntityNotFoundError("user", fmt.Sprintf("guid:%s", guid))
	}
	if err := s.userRepository.Delete(guid); err != nil {
		return exception.InternalError("failed user delete")
	}
	return nil
}

func (s *UserService) List(limit, offset int) ([]entity.User, *exception.ApiError) {
	users, err := s.userRepository.List(limit, offset)
	if err != nil {
		return nil, exception.InternalError("failed list users")
	}
	return users, nil
}

func (s *UserService) GetByGuid(guid uuid.UUID) (*entity.User, *exception.ApiError) {
	user, err := s.userRepository.FindByGuid(guid)
	if err != nil {
		return nil, exception.EntityNotFoundError("user", fmt.Sprintf("guid:%s", guid))
	}
	return user, nil
}

func (s *UserService) EnsureAdmin(cfg config.Config) *exception.ApiError {
	_, err := s.userRepository.FindByLogin(cfg.AdminLogin)
	if err == nil {
		return nil
	}

	hashedPassword, err := HashPassword(cfg.AdminPassword)
	if err != nil {
		return exception.InternalError("failed password hashing: create admin")
	}

	_, err = s.userRepository.Create(entity.User{
		Login:     cfg.AdminLogin,
		Password:  hashedPassword,
		IsAdmin:   true,
		IsBlocked: false,
	})
	if err != nil {
		return exception.InternalError("failed admin user creating")
	}

	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
