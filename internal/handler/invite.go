package handler

import (
	"net/http"
	"time"

	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/service"

	"github.com/gin-gonic/gin"
)

type IInviteHandler interface {
	Create(c *gin.Context)
	ListByTeam(c *gin.Context)
	Delete(c *gin.Context)
	Get(c *gin.Context)
}

type InviteHandler struct {
	BaseHandler
	services *service.Service
}

func NewInviteHandler(services *service.Service) *InviteHandler {
	return &InviteHandler{
		BaseHandler: *NewBaseHandler(nil),
		services:    services,
	}
}

func (h *InviteHandler) Create(c *gin.Context) {
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

	invite, apiErr := h.services.Invite.Create(teamGuid, user.Guid)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusCreated, dto.InviteCreateResponseDto{
		Guid:      invite.Guid,
		TeamGuid:  invite.TeamGuid,
		ExpiredAt: invite.ExpiredAt,
	})
}

func (h *InviteHandler) ListByTeam(c *gin.Context) {
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

	invites, apiErr := h.services.Invite.ListByTeam(teamGuid, user.Guid)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	response := make([]dto.InviteCreateResponseDto, 0, len(invites))
	for _, inv := range invites {
		response = append(response, dto.InviteCreateResponseDto{
			Guid:      inv.Guid,
			TeamGuid:  inv.TeamGuid,
			ExpiredAt: inv.ExpiredAt,
		})
	}
	c.JSON(http.StatusOK, response)
}

func (h *InviteHandler) Delete(c *gin.Context) {
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

	inviteGuid, ok := h.parseUUID(c, "invite_guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid invite guid"))
		return
	}

	if apiErr := h.services.Invite.DeleteByLeader(teamGuid, inviteGuid, user.Guid); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.EmptyResponse{})
}

func (h *InviteHandler) Get(c *gin.Context) {
	inviteGuid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid invite guid"))
		return
	}

	invite, apiErr := h.services.Invite.GetByGuid(inviteGuid)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.InviteCheckResponseDto{
		Guid:      invite.Guid,
		TeamGuid:  invite.TeamGuid,
		ExpiredAt: invite.ExpiredAt,
		IsExpired: time.Now().UTC().After(invite.ExpiredAt),
	})
}
