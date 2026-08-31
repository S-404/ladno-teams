package handler

import (
	"net/http"

	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IUserHandler interface {
	Me(c *gin.Context)
	UpdateMe(c *gin.Context)
	UpdateProfile(c *gin.Context)
}

type UserHandler struct {
	BaseHandler
	services *service.Service
}

func NewUserHandler(services *service.Service, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		BaseHandler: *NewBaseHandler(validate),
		services:    services,
	}
}

func (h *UserHandler) Me(c *gin.Context) {
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	fullUser, apiErr := h.services.User.GetByLogin(user.Login)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	profile, _ := h.services.Profile.GetByUserGuid(fullUser.Guid)

	c.JSON(http.StatusOK, gin.H{
		"user":    fullUser,
		"profile": profile,
	})
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	var req dto.UserUpdateRequestDto
	if err := h.ValidateRequestBody(c, &req); err != nil {
		exception.HttpResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	updated, apiErr := h.services.User.Update(user.Guid, req)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	var req dto.ProfileUpdateRequestDto
	if err := h.ValidateRequestBody(c, &req); err != nil {
		exception.HttpResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	profile, apiErr := h.services.Profile.Update(user.Guid, req)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, profile)
}
