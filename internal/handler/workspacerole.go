package handler

import (
	"net/http"

	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IWorkspaceRoleHandler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type WorkspaceRoleHandler struct {
	BaseHandler
	services *service.Service
}

func NewWorkspaceRoleHandler(services *service.Service, validate *validator.Validate) *WorkspaceRoleHandler {
	return &WorkspaceRoleHandler{
		BaseHandler: *NewBaseHandler(validate),
		services:    services,
	}
}

func (h *WorkspaceRoleHandler) Create(c *gin.Context) {
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	workspaceGuid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid workspace guid"))
		return
	}

	var req dto.WorkspaceRoleCreateRequestDto
	if err := h.ValidateRequestBody(c, &req); err != nil {
		exception.HttpResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	role, apiErr := h.services.WorkspaceRole.Create(user.Guid, workspaceGuid, req)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusCreated, role)
}

func (h *WorkspaceRoleHandler) List(c *gin.Context) {
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	workspaceGuid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid workspace guid"))
		return
	}

	roles, apiErr := h.services.WorkspaceRole.ListByWorkspace(user.Guid, workspaceGuid)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	response := make([]dto.WorkspaceRoleResponseDto, 0, len(roles))
	for _, r := range roles {
		response = append(response, dto.WorkspaceRoleResponseDto{
			WorkspaceGuid: r.WorkspaceGuid,
			TeammateGuid:  r.TeammateGuid,
			Role:          r.Role,
			CreatedAt:     r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *WorkspaceRoleHandler) Update(c *gin.Context) {
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	workspaceGuid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid workspace guid"))
		return
	}

	teammateGuid, ok := h.parseUUID(c, "teammate_guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid teammate guid"))
		return
	}

	var req dto.WorkspaceRoleUpdateRequestDto
	if err := h.ValidateRequestBody(c, &req); err != nil {
		exception.HttpResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	role, apiErr := h.services.WorkspaceRole.Update(user.Guid, workspaceGuid, teammateGuid, req)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, role)
}

func (h *WorkspaceRoleHandler) Delete(c *gin.Context) {
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	workspaceGuid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid workspace guid"))
		return
	}

	teammateGuid, ok := h.parseUUID(c, "teammate_guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid teammate guid"))
		return
	}

	if apiErr := h.services.WorkspaceRole.Delete(user.Guid, workspaceGuid, teammateGuid); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.EmptyResponse{})
}
