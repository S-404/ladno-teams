package handler

import (
	"net/http"
	"strconv"

	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ITeamHandler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type TeamHandler struct {
	BaseHandler
	services *service.Service
}

func NewTeamHandler(services *service.Service, validate *validator.Validate) *TeamHandler {
	return &TeamHandler{
		BaseHandler: *NewBaseHandler(validate),
		services:    services,
	}
}

func (h *TeamHandler) Create(c *gin.Context) {
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	var req dto.TeamCreateRequestDto
	if err := h.ValidateRequestBody(c, &req); err != nil {
		exception.HttpResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	team, apiErr := h.services.Team.Create(user.Guid, req)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusCreated, team)
}

func (h *TeamHandler) List(c *gin.Context) {
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	limit := h.parseLimit(c)
	offset := h.parseOffset(c)

	teams, apiErr := h.services.Team.ListByUser(user.Guid, limit, offset)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, teams)
}

func (h *TeamHandler) Get(c *gin.Context) {
	teamGuid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid team guid"))
		return
	}

	team, apiErr := h.services.Team.GetByGuid(teamGuid)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) Update(c *gin.Context) {
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

	var req dto.TeamUpdateRequestDto
	if err := h.ValidateRequestBody(c, &req); err != nil {
		exception.HttpResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	team, apiErr := h.services.Team.Update(user.Guid, teamGuid, req)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) Delete(c *gin.Context) {
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

	if apiErr := h.services.Team.Delete(user.Guid, teamGuid); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.EmptyResponse{})
}

func (h *TeamHandler) parseLimit(c *gin.Context) int {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil || limit <= 0 || limit > 100 {
		return 50
	}
	return limit
}

func (h *TeamHandler) parseOffset(c *gin.Context) int {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		return 0
	}
	return offset
}

func (h *TeamHandler) requireUser(c *gin.Context) uuid.UUID {
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return uuid.Nil
	}
	return user.Guid
}
