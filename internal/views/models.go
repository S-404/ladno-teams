package views

import (
	"html/template"

	"ladno-teams/internal/entity"
)

const AdminModalFormHTMX = `hx-target="#content" hx-swap="innerHTML" hx-on::after-request="if(event.detail.successful) closeAdminModal()"`

type LayoutData struct {
	Title   string
	Content template.HTML
}

type SectionHeader struct {
	Title       string
	ActionLabel string
	ActionURL   string
}

type CheckboxForm struct {
	PostURL   string
	FieldName string
	Value     string
	Checked   bool
}

type AdminUsersSection struct {
	Users []entity.AdminUserView
}

type AdminInvitesSection struct {
	Invites []entity.AdminInviteView
}

type AdminTeamsSection struct {
	Teams []entity.Team
}

type AdminTeammatesSection struct {
	Teammates []entity.AdminTeammateView
}

type AdminWorkspacesSection struct {
	Workspaces []entity.AdminWorkspaceView
}

type AdminWorkspaceRolesSection struct {
	Roles []entity.AdminWorkspaceRoleView
}

type AdminInviteModal struct {
	Teams            []entity.Team
	SelectedTeamGUID string
}

type TeamSelect struct {
	Teams    []entity.Team
	Selected string
}

type AdminWorkspaceModal struct {
	Teams []entity.Team
}

type AdminWorkspaceRoleModal struct {
	Workspaces []entity.AdminWorkspaceView
	Users      []entity.AdminUserView
}

type AuthPage struct {
	Title        string
	ErrorMessage string
	Login        string
}

type AuthRegisterPage struct {
	Title        string
	ErrorMessage string
	GUID         string
	TeamName     string
	IsExpired    bool
}

type AuthHomePage struct {
	Title string
	Login string
}
