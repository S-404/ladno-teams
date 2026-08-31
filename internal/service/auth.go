package service

import (
	"errors"
	"ladno-teams/internal/entity"
	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type IAuthService interface {
	Auth(user entity.User) (*dto.AuthTokens, *exception.ApiError)
	RefreshAccessToken(refreshToken string) (*dto.AuthTokens, *exception.ApiError)
	AuthCookie(refreshToken string) http.Cookie
	ParseToken(accessToken string) (*dto.TokenPayload, error)
}

type AuthService struct {
	BaseService
	sessionService ISessionService
}

func NewAuthService(baseService *BaseService, sessionService ISessionService) *AuthService {
	return &AuthService{
		BaseService:    *baseService,
		sessionService: sessionService,
	}
}

func (s *AuthService) Auth(user entity.User) (*dto.AuthTokens, *exception.ApiError) {
	return s.generateAuthTokens(user)
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (*dto.AuthTokens, *exception.ApiError) {
	parsedTokenPayload, err := s.ParseToken(refreshToken)
	if err != nil {
		return nil, exception.AuthError("Invalid Token")
	}

	if _, apiErr := s.sessionService.FindByToken(refreshToken); apiErr != nil {
		return nil, apiErr
	}

	if apiErr := s.sessionService.DestroyByToken(refreshToken); apiErr != nil {
		return nil, apiErr
	}

	user := entity.User{
		Guid:    parsedTokenPayload.User.Guid,
		Login:   parsedTokenPayload.User.Login,
		IsAdmin: parsedTokenPayload.User.IsAdmin,
	}

	return s.generateAuthTokens(user)
}

func (s *AuthService) AuthCookie(refreshToken string) http.Cookie {
	return http.Cookie{
		Name:     "token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.cfg.CookieSecure,
		MaxAge:   s.cfg.CookieMaxAge,
	}
}

func (s *AuthService) ParseToken(accessToken string) (*dto.TokenPayload, error) {
	token, err := jwt.ParseWithClaims(accessToken, &dto.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.cfg.JwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*dto.TokenClaims)
	if !ok {
		return nil, errors.New("token claims are not of type *dto.TokenClaims")
	}

	return &claims.TokenPayload, nil
}

func (s *AuthService) generateAuthTokens(user entity.User) (*dto.AuthTokens, *exception.ApiError) {
	accessToken, err := s.generateToken(user, s.cfg.JwtAccessTTL)
	if err != nil {
		return nil, exception.InternalError("failed generate access token")
	}

	refreshToken, err := s.generateToken(user, s.cfg.JwtRefreshTTL)
	if err != nil {
		return nil, exception.InternalError("failed generate refresh token")
	}

	if _, apiErr := s.sessionService.Create(user.Guid, refreshToken, time.Duration(s.cfg.JwtRefreshTTL)*time.Second); apiErr != nil {
		return nil, apiErr
	}

	return &dto.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateToken(user entity.User, duration int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &dto.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   "user",
			ExpiresAt: time.Now().Add(time.Second * time.Duration(duration)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		TokenPayload: dto.TokenPayload{
			User: dto.UserDto{
				Guid:      user.Guid,
				Login:     user.Login,
				IsAdmin:   user.IsAdmin,
				IsBlocked: user.IsBlocked,
			},
		},
	})

	return token.SignedString([]byte(s.cfg.JwtSecretKey))
}
