package handler

import (
	"ladno-teams/internal/entity"
	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IAdminHandler interface {
	Index(c *gin.Context)
	Users(c *gin.Context)
	Invites(c *gin.Context)
	Teams(c *gin.Context)
	Teammates(c *gin.Context)
	Workspaces(c *gin.Context)
	WorkspaceRoles(c *gin.Context)
	ModalInvite(c *gin.Context)
	ModalTeam(c *gin.Context)
	ModalWorkspaceRole(c *gin.Context)
	CreateInvite(c *gin.Context)
	CreateTeam(c *gin.Context)
	CreateWorkspaceRole(c *gin.Context)
	ToggleUserBlock(c *gin.Context)
	ToggleUserAdmin(c *gin.Context)
	DeleteUser(c *gin.Context)
	DeleteTeam(c *gin.Context)
	DeleteTeammate(c *gin.Context)
	DeleteInvite(c *gin.Context)
	DeleteWorkspace(c *gin.Context)
	DeleteWorkspaceRole(c *gin.Context)
}

type AdminHandler struct {
	BaseHandler
	services *service.Service
}

func NewAdminHandler(services *service.Service, validate *validator.Validate) *AdminHandler {
	return &AdminHandler{
		BaseHandler: *NewBaseHandler(validate),
		services:    services,
	}
}

func (h *AdminHandler) Index(c *gin.Context) {
	renderHTML(c, adminLayout(adminDashboard()))
}

func (h *AdminHandler) Users(c *gin.Context) {
	users, apiErr := h.services.Admin.ListUsers(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTML(c, adminUsersSection(users))
}

func (h *AdminHandler) Invites(c *gin.Context) {
	invites, apiErr := h.services.Admin.ListInvites(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTML(c, adminInvitesSection(invites))
}

func (h *AdminHandler) Teams(c *gin.Context) {
	teams, apiErr := h.services.Team.ListAll(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTML(c, adminTeamsSection(teams))
}

func (h *AdminHandler) Teammates(c *gin.Context) {
	teammates, apiErr := h.services.Admin.ListTeammates(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTML(c, adminTeammatesSection(teammates))
}

func (h *AdminHandler) Workspaces(c *gin.Context) {
	workspaces, apiErr := h.services.Admin.ListWorkspaces(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTML(c, adminWorkspacesSection(workspaces))
}

func (h *AdminHandler) WorkspaceRoles(c *gin.Context) {
	roles, apiErr := h.services.Admin.ListWorkspaceRoles(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTML(c, adminWorkspaceRolesSection(roles))
}

func (h *AdminHandler) ModalInvite(c *gin.Context) {
	teams, apiErr := h.services.Team.ListAll(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTML(c, adminInviteModal(teams, c.Query("team_guid")))
}

func (h *AdminHandler) ModalTeam(c *gin.Context) {
	renderHTML(c, adminTeamModal())
}

func (h *AdminHandler) ModalWorkspaceRole(c *gin.Context) {
	workspaces, apiErr := h.services.Admin.ListWorkspaces(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	users, apiErr := h.services.Admin.ListUsers(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTML(c, adminWorkspaceRoleModal(workspaces, users))
}

func (h *AdminHandler) CreateInvite(c *gin.Context) {
	teamGuid, ok := h.parseUUIDFromForm(c, "team_guid")
	if !ok {
		renderHTML(c, adminInviteModal(nil, "")+"<p>Invalid team</p>")
		return
	}

	if _, apiErr := h.services.Invite.AdminCreate(teamGuid); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	invites, apiErr := h.services.Admin.ListInvites(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTMLCloseModal(c, adminInvitesSection(invites))
}

func (h *AdminHandler) CreateTeam(c *gin.Context) {
	adminUser, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")
	var descPtr *string
	if description != "" {
		descPtr = &description
	}

	if _, apiErr := h.services.Team.Create(adminUser.Guid, dto.TeamCreateRequestDto{
		Name:        name,
		Description: descPtr,
	}); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	teams, apiErr := h.services.Team.ListAll(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTMLCloseModal(c, adminTeamsSection(teams))
}

func (h *AdminHandler) CreateWorkspaceRole(c *gin.Context) {
	workspaceGuid, ok := h.parseUUIDFromForm(c, "workspace_guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid workspace guid"))
		return
	}
	teammateGuid, ok := h.parseUUIDFromForm(c, "teammate_guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid teammate guid"))
		return
	}

	role := entity.Role(c.PostForm("role"))
	if _, apiErr := h.services.WorkspaceRole.AdminCreate(workspaceGuid, dto.WorkspaceRoleCreateRequestDto{
		TeammateGuid: teammateGuid,
		Role:         role,
	}); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	roles, apiErr := h.services.Admin.ListWorkspaceRoles(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTMLCloseModal(c, adminWorkspaceRolesSection(roles))
}

func (h *AdminHandler) ToggleUserBlock(c *gin.Context) {
	guid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid user guid"))
		return
	}

	user, apiErr := h.services.User.GetByGuid(guid)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	isBlocked := !user.IsBlocked
	if _, apiErr := h.services.User.AdminUpdate(guid, dto.UserAdminUpdateRequestDto{IsBlocked: &isBlocked}); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	view, apiErr := h.services.Admin.GetUserViewByGuid(guid.String())
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTML(c, adminUserRow(*view))
}

func (h *AdminHandler) ToggleUserAdmin(c *gin.Context) {
	guid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid user guid"))
		return
	}

	user, apiErr := h.services.User.GetByGuid(guid)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	isAdmin := !user.IsAdmin
	if _, apiErr := h.services.User.AdminUpdate(guid, dto.UserAdminUpdateRequestDto{IsAdmin: &isAdmin}); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	view, apiErr := h.services.Admin.GetUserViewByGuid(guid.String())
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTML(c, adminUserRow(*view))
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	guid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid user guid"))
		return
	}

	if apiErr := h.services.User.Delete(guid); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AdminHandler) DeleteTeam(c *gin.Context) {
	guid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid team guid"))
		return
	}

	if apiErr := h.services.Team.AdminDelete(guid); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AdminHandler) DeleteTeammate(c *gin.Context) {
	userGuid, ok := h.parseUUID(c, "user_guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid user guid"))
		return
	}
	teamGuid, ok := h.parseUUID(c, "team_guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid team guid"))
		return
	}

	if apiErr := h.services.Teammate.AdminDelete(userGuid, teamGuid); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AdminHandler) DeleteInvite(c *gin.Context) {
	guid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid invite guid"))
		return
	}

	if apiErr := h.services.Invite.Delete(guid); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AdminHandler) DeleteWorkspace(c *gin.Context) {
	guid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid workspace guid"))
		return
	}

	if apiErr := h.services.Workspace.AdminDelete(guid); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AdminHandler) DeleteWorkspaceRole(c *gin.Context) {
	workspaceGuid, ok := h.parseUUID(c, "workspace_guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid workspace guid"))
		return
	}
	teammateGuid, ok := h.parseUUID(c, "teammate_guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid teammate guid"))
		return
	}

	if apiErr := h.services.WorkspaceRole.AdminDelete(workspaceGuid, teammateGuid); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AdminHandler) parseUUIDFromForm(c *gin.Context, field string) (uuid.UUID, bool) {
	value, err := uuid.Parse(c.PostForm(field))
	if err != nil {
		return uuid.Nil, false
	}
	return value, true
}

func renderHTML(c *gin.Context, content string) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
}

func renderHTMLCloseModal(c *gin.Context, content string) {
	c.Header("HX-Trigger", `{"closeModal":{"target":"body"}}`)
	renderHTML(c, content+`<div id="modal-host" hx-swap-oob="true"></div>`)
}
