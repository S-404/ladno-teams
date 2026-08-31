package handler

import (
	"fmt"
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
	Profiles(c *gin.Context)
	Teams(c *gin.Context)
	Teammates(c *gin.Context)
	Invites(c *gin.Context)
	Workspaces(c *gin.Context)
	WorkspaceRoles(c *gin.Context)
	ToggleUserBlock(c *gin.Context)
	ToggleUserAdmin(c *gin.Context)
	DeleteUser(c *gin.Context)
	DeleteProfile(c *gin.Context)
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
	users, apiErr := h.services.User.List(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTML(c, adminUsersTable(users))
}

func (h *AdminHandler) Profiles(c *gin.Context) {
	users, apiErr := h.services.User.List(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	profiles := make([]entity.Profile, 0, len(users))
	for _, u := range users {
		if p, err := h.services.Profile.GetByUserGuid(u.Guid); err == nil {
			profiles = append(profiles, *p)
		}
	}
	renderHTML(c, adminProfilesTable(profiles))
}

func (h *AdminHandler) Teams(c *gin.Context) {
	teams, apiErr := h.services.Team.ListAll(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTML(c, adminTeamsTable(teams))
}

func (h *AdminHandler) Teammates(c *gin.Context) {
	teams, apiErr := h.services.Team.ListAll(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	var teammates []entity.Teammate
	for _, t := range teams {
		list, apiErr := h.services.Teammate.ListByTeam(t.Guid)
		if apiErr == nil {
			teammates = append(teammates, list...)
		}
	}
	renderHTML(c, adminTeammatesTable(teammates))
}

func (h *AdminHandler) Invites(c *gin.Context) {
	teams, apiErr := h.services.Team.ListAll(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	var invites []entity.Invite
	for _, t := range teams {
		// invites are consumed on use, so for admin we could query all invites. We'll skip listing by team for now.
		_ = t
	}
	renderHTML(c, adminInvitesTable(invites))
}

func (h *AdminHandler) Workspaces(c *gin.Context) {
	workspaces, apiErr := h.services.Workspace.ListAll(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	renderHTML(c, adminWorkspacesTable(workspaces))
}

func (h *AdminHandler) WorkspaceRoles(c *gin.Context) {
	workspaces, apiErr := h.services.Workspace.ListAll(1000, 0)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}
	var roles []entity.WorkspaceRole
	for _, w := range workspaces {
		list, apiErr := h.services.WorkspaceRole.ListByWorkspace(uuid.Nil, w.Guid)
		if apiErr == nil {
			roles = append(roles, list...)
		}
	}
	renderHTML(c, adminWorkspaceRolesTable(roles))
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
	updated, apiErr := h.services.User.AdminUpdate(guid, dto.UserAdminUpdateRequestDto{IsBlocked: &isBlocked})
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	renderHTML(c, adminUserRow(*updated))
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
	updated, apiErr := h.services.User.AdminUpdate(guid, dto.UserAdminUpdateRequestDto{IsAdmin: &isAdmin})
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	renderHTML(c, adminUserRow(*updated))
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

func (h *AdminHandler) DeleteProfile(c *gin.Context) {
	guid, ok := h.parseUUID(c, "guid")
	if !ok {
		exception.HttpResponseException(c, exception.BadRequest("invalid user guid"))
		return
	}

	if apiErr := h.services.Profile.DeleteByUserGuid(guid); apiErr != nil {
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

func renderHTML(c *gin.Context, content string) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
}

func adminLayout(content string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Admin Panel</title>
	<script src="https://unpkg.com/htmx.org@1.9.12"></script>
	<style>
		body { font-family: sans-serif; margin: 0; background: #f5f5f5; }
		.navbar { background: #333; padding: 1rem; color: white; display: flex; gap: 1rem; align-items: center; }
		.navbar a { color: white; text-decoration: none; cursor: pointer; }
		.navbar a:hover { text-decoration: underline; }
		.container { padding: 1rem; }
		table { width: 100%%; border-collapse: collapse; background: white; }
		th, td { padding: 0.5rem; border: 1px solid #ddd; text-align: left; }
		th { background: #eee; }
		button { cursor: pointer; }
		tr.htmx-swapping { opacity: 0; transition: opacity 0.3s; }
	</style>
</head>
<body>
	<div class="navbar">
		<a hx-get="/admin/users" hx-target="#content" hx-push-url="false">Users</a>
		<a hx-get="/admin/profiles" hx-target="#content" hx-push-url="false">Profiles</a>
		<a hx-get="/admin/teams" hx-target="#content" hx-push-url="false">Teams</a>
		<a hx-get="/admin/teammates" hx-target="#content" hx-push-url="false">Teammates</a>
		<a hx-get="/admin/invites" hx-target="#content" hx-push-url="false">Invites</a>
		<a hx-get="/admin/workspaces" hx-target="#content" hx-push-url="false">Workspaces</a>
		<a hx-get="/admin/workspace_roles" hx-target="#content" hx-push-url="false">Workspace Roles</a>
		<form method="POST" action="/admin/logout" style="margin-left:auto;">
			<button type="submit" style="background:#555;color:white;border:none;padding:0.4rem 0.8rem;border-radius:4px;">Logout</button>
		</form>
	</div>
	<div class="container" id="content">
		%s
	</div>
</body>
</html>`, content)
}

func adminDashboard() string {
	return "<h1>Admin Dashboard</h1><p>Выберите раздел в меню.</p>"
}

func adminUsersTable(users []entity.User) string {
	rows := ""
	for _, u := range users {
		rows += adminUserRow(u)
	}
	return fmt.Sprintf(`<h1>Users</h1>
<table>
<thead><tr><th>Guid</th><th>Login</th><th>Admin</th><th>Blocked</th><th>Actions</th></tr></thead>
<tbody>%s</tbody>
</table>`, rows)
}

func adminUserRow(u entity.User) string {
	adminLabel := "Make admin"
	if u.IsAdmin {
		adminLabel = "Revoke admin"
	}
	blockLabel := "Block"
	if u.IsBlocked {
		blockLabel = "Unblock"
	}
	return fmt.Sprintf(`<tr id="user-%s">
	<td>%s</td>
	<td>%s</td>
	<td>%v</td>
	<td>%v</td>
	<td>
		<button hx-post="/admin/users/%s/toggle-admin" hx-target="#user-%s" hx-swap="outerHTML">%s</button>
		<button hx-post="/admin/users/%s/toggle-block" hx-target="#user-%s" hx-swap="outerHTML">%s</button>
		<button hx-delete="/admin/users/%s" hx-target="#user-%s" hx-swap="outerHTML swap:0.3s" hx-confirm="Delete user?">Delete</button>
	</td>
</tr>`, u.Guid, u.Guid, u.Login, u.IsAdmin, u.IsBlocked, u.Guid, u.Guid, adminLabel, u.Guid, u.Guid, blockLabel, u.Guid, u.Guid)
}

func adminProfilesTable(profiles []entity.Profile) string {
	rows := ""
	for _, p := range profiles {
		about := ""
		if p.About != nil {
			about = *p.About
		}
		rows += fmt.Sprintf(`<tr id="profile-%s">
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td><button hx-delete="/admin/profiles/%s" hx-target="#profile-%s" hx-swap="outerHTML swap:0.3s" hx-confirm="Delete profile?">Delete</button></td>
</tr>`, p.UserGuid, p.UserGuid, p.Name, about, p.UserGuid, p.UserGuid)
	}
	return fmt.Sprintf(`<h1>Profiles</h1>
<table>
<thead><tr><th>User Guid</th><th>Name</th><th>About</th><th>Actions</th></tr></thead>
<tbody>%s</tbody>
</table>`, rows)
}

func adminTeamsTable(teams []entity.Team) string {
	rows := ""
	for _, t := range teams {
		desc := ""
		if t.Description != nil {
			desc = *t.Description
		}
		rows += fmt.Sprintf(`<tr id="team-%s">
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td><button hx-delete="/admin/teams/%s" hx-target="#team-%s" hx-swap="outerHTML swap:0.3s" hx-confirm="Delete team?">Delete</button></td>
</tr>`, t.Guid, t.Guid, t.Name, desc, t.Guid, t.Guid)
	}
	return fmt.Sprintf(`<h1>Teams</h1>
<table>
<thead><tr><th>Guid</th><th>Name</th><th>Description</th><th>Actions</th></tr></thead>
<tbody>%s</tbody>
</table>`, rows)
}

func adminTeammatesTable(teammates []entity.Teammate) string {
	rows := ""
	for _, tm := range teammates {
		rows += fmt.Sprintf(`<tr id="teammate-%s-%s">
	<td>%s</td>
	<td>%s</td>
	<td>%v</td>
	<td><button hx-delete="/admin/teammates/%s/%s" hx-target="#teammate-%s-%s" hx-swap="outerHTML swap:0.3s" hx-confirm="Delete teammate?">Delete</button></td>
</tr>`, tm.UserGuid, tm.TeamGuid, tm.UserGuid, tm.TeamGuid, tm.IsLeader, tm.UserGuid, tm.TeamGuid, tm.UserGuid, tm.TeamGuid)
	}
	return fmt.Sprintf(`<h1>Teammates</h1>
<table>
<thead><tr><th>User Guid</th><th>Team Guid</th><th>Leader</th><th>Actions</th></tr></thead>
<tbody>%s</tbody>
</table>`, rows)
}

func adminInvitesTable(invites []entity.Invite) string {
	rows := ""
	for _, inv := range invites {
		rows += fmt.Sprintf(`<tr id="invite-%s">
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td><button hx-delete="/admin/invites/%s" hx-target="#invite-%s" hx-swap="outerHTML swap:0.3s" hx-confirm="Delete invite?">Delete</button></td>
</tr>`, inv.Guid, inv.Guid, inv.TeamGuid, inv.ExpiredAt.Format("2006-01-02 15:04:05"), inv.Guid, inv.Guid)
	}
	return fmt.Sprintf(`<h1>Invites</h1>
<table>
<thead><tr><th>Guid</th><th>Team Guid</th><th>Expired At</th><th>Actions</th></tr></thead>
<tbody>%s</tbody>
</table>`, rows)
}

func adminWorkspacesTable(workspaces []entity.Workspace) string {
	rows := ""
	for _, w := range workspaces {
		rows += fmt.Sprintf(`<tr id="workspace-%s">
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td><button hx-delete="/admin/workspaces/%s" hx-target="#workspace-%s" hx-swap="outerHTML swap:0.3s" hx-confirm="Delete workspace?">Delete</button></td>
</tr>`, w.Guid, w.Guid, w.Name, w.Version, w.TeamGuid, w.Guid, w.Guid)
	}
	return fmt.Sprintf(`<h1>Workspaces</h1>
<table>
<thead><tr><th>Guid</th><th>Name</th><th>Version</th><th>Team Guid</th><th>Actions</th></tr></thead>
<tbody>%s</tbody>
</table>`, rows)
}

func adminWorkspaceRolesTable(roles []entity.WorkspaceRole) string {
	rows := ""
	for _, r := range roles {
		rows += fmt.Sprintf(`<tr id="wr-%s-%s">
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td><button hx-delete="/admin/workspace_roles/%s/%s" hx-target="#wr-%s-%s" hx-swap="outerHTML swap:0.3s" hx-confirm="Delete workspace role?">Delete</button></td>
</tr>`, r.WorkspaceGuid, r.TeammateGuid, r.WorkspaceGuid, r.TeammateGuid, r.Role, r.WorkspaceGuid, r.TeammateGuid, r.WorkspaceGuid, r.TeammateGuid)
	}
	return fmt.Sprintf(`<h1>Workspace Roles</h1>
<table>
<thead><tr><th>Workspace Guid</th><th>Teammate Guid</th><th>Role</th><th>Actions</th></tr></thead>
<tbody>%s</tbody>
</table>`, rows)
}
