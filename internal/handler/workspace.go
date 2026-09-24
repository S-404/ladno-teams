package handler

import (
	"net/http"
	"strconv"

	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IWorkspaceHandler interface {
	Create(c *gin.Context)
	ListByTeam(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type WorkspaceHandler struct {
	BaseHandler
	services *service.Service
}

func NewWorkspaceHandler(services *service.Service, validate *validator.Validate) *WorkspaceHandler {
	return &WorkspaceHandler{
		BaseHandler: *NewBaseHandler(validate),
		services:    services,
	}
}

func (h *WorkspaceHandler) Create(c *gin.Context) {
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

	var req dto.WorkspaceCreateRequestDto
	if err := h.ValidateRequestBody(c, &req); err != nil {
		exception.HttpResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	workspace, apiErr := h.services.Workspace.Create(user.Guid, teamGuid, req)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusCreated, workspace)
}

func (h *WorkspaceHandler) ListByTeam(c *gin.Context) {
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

	limit := h.parseLimit(c)
	offset := h.parseOffset(c)

	workspaces, apiErr := h.services.Workspace.ListByTeam(user.Guid, teamGuid, limit, offset)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	response := make([]dto.WorkspaceListResponseDto, 0, len(workspaces))
	for _, w := range workspaces {
		response = append(response, dto.WorkspaceListResponseDto{
			Guid:      w.Guid,
			TeamGuid:  w.TeamGuid,
			Name:      w.Name,
			Version:   w.Version,
			GitURL:    "/git/workspaces/" + w.Guid.String(),
			CreatedAt: w.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: w.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *WorkspaceHandler) Get(c *gin.Context) {
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

	workspace, apiErr := h.services.Workspace.GetByGuid(user.Guid, workspaceGuid)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, workspace)
}

func (h *WorkspaceHandler) Update(c *gin.Context) {
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

	var req dto.WorkspaceUpdateRequestDto
	if err := h.ValidateRequestBody(c, &req); err != nil {
		exception.HttpResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	workspace, apiErr := h.services.Workspace.Update(user.Guid, workspaceGuid, req)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, workspace)
}

func (h *WorkspaceHandler) Delete(c *gin.Context) {
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

	if apiErr := h.services.Workspace.Delete(user.Guid, workspaceGuid); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.EmptyResponse{})
}

func (h *WorkspaceHandler) parseLimit(c *gin.Context) int {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil || limit <= 0 || limit > 100 {
		return 50
	}
	return limit
}

func (h *WorkspaceHandler) parseOffset(c *gin.Context) int {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		return 0
	}
	return offset
}
