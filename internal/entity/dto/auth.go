package dto

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type AuthRequestDto struct {
	Login    string `json:"login" binding:"required,valid_login"`
	Password string `json:"password" binding:"required,valid_password"`
}

type AuthResponseDto struct {
	AccessToken string `json:"accessToken"`
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

type TokenClaims struct {
	jwt.StandardClaims
	TokenPayload
}

type TokenPayload struct {
	User UserDto `json:"user"`
}

type InviteRegisterRequestDto struct {
	Guid     uuid.UUID `json:"guid" binding:"required"`
	Login    string    `json:"login" binding:"required,valid_login"`
	Password string    `json:"password" binding:"required,valid_password"`
	Name     string    `json:"name" binding:"required,max=255"`
}

type AcceptInviteRequestDto struct {
	Guid uuid.UUID `json:"guid" binding:"required"`
}
