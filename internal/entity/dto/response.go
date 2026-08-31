package dto

import "github.com/google/uuid"

type EmptyResponse struct{}

type ErrorResponse struct {
	Error ApiErrorDto `json:"error"`
}

type ApiErrorDto struct {
	Guid    string `json:"guid"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

type UserDto struct {
	Guid      uuid.UUID `json:"guid" format:"uuid"`
	Login     string    `json:"login"`
	IsAdmin   bool      `json:"is_admin"`
	IsBlocked bool      `json:"is_blocked"`
}

type ProfileDto struct {
	UserGuid uuid.UUID `json:"user_guid" format:"uuid"`
	Name     string    `json:"name"`
	About    *string   `json:"about,omitempty"`
}
