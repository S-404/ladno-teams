package handler

import (
	"net/http"

	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ITeammateHandler interface {
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type TeammateHandler struct {
	BaseHandler
	services *service.Service
}

func NewTeammateHandler(services *service.Service, validate *validator.Validate) *TeammateHandler {
	return &TeammateHandler{
		BaseHandler: *NewBaseHandler(validate),
		services:    services,
	}
}

func (h *TeammateHandler) List(c *gin.Context) {
	teamGuid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid team guid"))
		return
	}

	teammates, apiErr := h.services.Teammate.ListByTeam(teamGuid)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	response := make([]dto.TeammateListResponseDto, 0, len(teammates))
	for _, t := range teammates {
		user := dto.TeammateUserDto{
			Guid:  t.UserGuid,
			Login: t.UserLogin,
		}
		if t.ProfileName != nil && *t.ProfileName != "" {
			user.Profile = &dto.TeammateProfileDto{Name: *t.ProfileName}
		}
		response = append(response, dto.TeammateListResponseDto{
			UserGuid:  t.UserGuid,
			TeamGuid:  t.TeamGuid,
			IsLeader:  t.IsLeader,
			CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
			User:      user,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *TeammateHandler) Update(c *gin.Context) {
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	teamGuid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid team guid"))
		return
	}

	targetUserGuid, ok := h.parseUUID(c, "user_guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid user guid"))
		return
	}

	var req dto.TeammateUpdateRequestDto
	if err := h.ValidateRequestBody(c, &req); err != nil {
		exception.HttpResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	teammate, apiErr := h.services.Teammate.Update(user.Guid, teamGuid, targetUserGuid, req)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, teammate)
}

func (h *TeammateHandler) Delete(c *gin.Context) {
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	teamGuid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid team guid"))
		return
	}

	targetUserGuid, ok := h.parseUUID(c, "user_guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid user guid"))
		return
	}

	if apiErr := h.services.Teammate.Delete(user.Guid, teamGuid, targetUserGuid); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.EmptyResponse{})
}
