package handler

import (
	"fmt"
	"html"
	"ladno-teams/internal/entity"
)

const adminDateFormat = "2006-01-02 15:04:05"

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
		.navbar { background: #333; padding: 1rem; color: white; display: flex; gap: 1rem; align-items: center; flex-wrap: wrap; }
		.navbar a { color: white; text-decoration: none; }
		.navbar a:hover { text-decoration: underline; }
		.container { padding: 1rem; }
		.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; gap: 1rem; }
		.section-header h1 { margin: 0; font-size: 1.5rem; }
		table { width: 100%%; border-collapse: collapse; background: white; }
		th, td { padding: 0.5rem; border: 1px solid #ddd; text-align: left; vertical-align: top; }
		th { background: #eee; }
		button, .btn { cursor: pointer; border: 1px solid #ccc; background: #fff; padding: 0.35rem 0.6rem; border-radius: 4px; }
		.btn-primary { background: #333; color: #fff; border-color: #333; }
		.btn-danger { background: #c0392b; color: #fff; border-color: #c0392b; }
		.actions { display: flex; flex-wrap: wrap; gap: 0.35rem; }
		tr.htmx-swapping { opacity: 0; transition: opacity 0.3s; }
		.modal-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,0.45); display: flex; align-items: center; justify-content: center; z-index: 1000; }
		.modal-card { background: #fff; padding: 1.5rem; border-radius: 8px; width: 100%%; max-width: 420px; box-shadow: 0 8px 24px rgba(0,0,0,0.2); }
		.modal-card h2 { margin-top: 0; }
		.modal-card label { display: block; margin-bottom: 0.25rem; font-weight: 600; }
		.modal-card input, .modal-card select, .modal-card textarea { width: 100%%; padding: 0.5rem; margin-bottom: 1rem; box-sizing: border-box; }
		.modal-actions { display: flex; gap: 0.5rem; justify-content: flex-end; }
		.guid-cell { font-family: monospace; font-size: 0.85rem; word-break: break-all; }
		.admin-section { margin-bottom: 1rem; }
		.section-search { margin-bottom: 1rem; }
		.section-search input { width: 100%%; max-width: 360px; padding: 0.5rem 0.65rem; border: 1px solid #ccc; border-radius: 6px; font-size: 1rem; box-sizing: border-box; }
		.section-search input:focus { outline: none; border-color: #333; box-shadow: 0 0 0 2px rgba(0,0,0,0.08); }
	</style>
	<script>
		function copyText(text) {
			navigator.clipboard.writeText(text);
		}
		function copyInviteLink(guid) {
			copyText(window.location.origin + '/register?guid=' + guid);
		}
		function filterAdminTable(input) {
			const section = input.closest('.admin-section');
			if (!section) return;
			const query = input.value.toLowerCase().trim();
			section.querySelectorAll('tbody tr').forEach(function(row) {
				row.style.display = row.textContent.toLowerCase().includes(query) ? '' : 'none';
			});
		}
		window.filterAdminTable = filterAdminTable;
		document.addEventListener('closeModal', function() {
			const host = document.getElementById('modal-host');
			if (host) host.innerHTML = '';
		});
		document.addEventListener('htmx:beforeSwap', function(evt) {
			const xhr = evt.detail.xhr;
			if (!xhr || xhr.status < 200 || xhr.status >= 300) return;
			if (xhr.responseText && xhr.responseText.trim() !== '') return;
			const target = evt.detail.target;
			if (target && target.tagName === 'TR') {
				evt.detail.shouldSwap = false;
				target.classList.add('htmx-swapping');
				setTimeout(function() { target.remove(); }, 300);
			}
		});
	</script>
</head>
<body>
	<div class="navbar">
		<a hx-get="/admin/users" hx-target="#content" hx-swap="innerHTML">Users</a>
		<a hx-get="/admin/invites" hx-target="#content" hx-swap="innerHTML">Invites</a>
		<a hx-get="/admin/teams" hx-target="#content" hx-swap="innerHTML">Teams</a>
		<a hx-get="/admin/teammates" hx-target="#content" hx-swap="innerHTML">Teammates</a>
		<a hx-get="/admin/workspaces" hx-target="#content" hx-swap="innerHTML">Workspaces</a>
		<a hx-get="/admin/workspace_roles" hx-target="#content" hx-swap="innerHTML">Workspace Roles</a>
		<form method="POST" action="/admin/logout" style="margin-left:auto;">
			<button type="submit" class="btn-primary">Logout</button>
		</form>
	</div>
	<div class="container" id="content">%s</div>
	<div id="modal-host"></div>
</body>
</html>`, content)
}

func adminDashboard() string {
	return "<h1>Admin Dashboard</h1><p>Выберите раздел в меню.</p>"
}

func adminSectionHeader(title, actionLabel, actionURL string) string {
	action := ""
	if actionLabel != "" && actionURL != "" {
		action = fmt.Sprintf(
			`<button type="button" class="btn-primary" hx-get="%s" hx-target="#modal-host" hx-swap="innerHTML">%s</button>`,
			html.EscapeString(actionURL),
			html.EscapeString(actionLabel),
		)
	}
	return fmt.Sprintf(`<div class="section-header"><h1>%s</h1>%s</div>`, html.EscapeString(title), action)
}

func adminSearchBar() string {
	return `<div class="section-search"><input type="search" class="admin-search-input" placeholder="Поиск..." autocomplete="off" oninput="filterAdminTable(this)"></div>`
}

func adminListSection(header, tableInner string) string {
	return fmt.Sprintf(`<div class="admin-section">%s%s<table class="admin-table">%s</table></div>`, header, adminSearchBar(), tableInner)
}

func adminBoolCheckboxForm(postURL, fieldName string, checked bool) string {
	checkedAttr := ""
	value := "false"
	if checked {
		checkedAttr = " checked"
		value = "true"
	}
	return fmt.Sprintf(`<form hx-post="%s" hx-trigger="change" hx-swap="none"
		hx-on::after-request="if(!event.detail.successful){const cb=this.querySelector('input[type=checkbox]');if(cb)cb.checked=!cb.checked;}">
		<input type="hidden" name="%s" value="%s">
		<input type="checkbox"%s onchange="this.form.querySelector('input[name=%s]').value=this.checked?'true':'false'">
	</form>`,
		html.EscapeString(postURL),
		html.EscapeString(fieldName),
		value,
		checkedAttr,
		html.EscapeString(fieldName),
	)
}

func adminUsersSection(users []entity.AdminUserView) string {
	rows := ""
	for _, u := range users {
		rows += adminUserRow(u)
	}
	return adminListSection(adminSectionHeader("Users", "Invite", "/admin/modals/invite"), fmt.Sprintf(`
<thead><tr><th>Login</th><th>Admin</th><th>Blocked</th><th>Name</th><th>About</th><th>Created</th><th>Updated</th><th>Actions</th></tr></thead>
<tbody>%s</tbody>`, rows))
}

func adminUserRow(u entity.AdminUserView) string {
	guid := u.Guid.String()
	return fmt.Sprintf(`<tr id="user-%s">
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td>%s</td><td>%s</td>
	<td>%s</td><td>%s</td>
	<td class="actions">
		<button class="btn-danger" hx-delete="/admin/users/%s" hx-target="closest tr" hx-swap="outerHTML settle:0.3s" hx-confirm="Delete user?">Delete</button>
	</td>
</tr>`,
		guid,
		html.EscapeString(u.Login),
		adminBoolCheckboxForm("/admin/users/"+guid+"/admin", "is_admin", u.IsAdmin),
		adminBoolCheckboxForm("/admin/users/"+guid+"/blocked", "is_blocked", u.IsBlocked),
		html.EscapeString(strVal(u.Name)), html.EscapeString(strVal(u.About)),
		html.EscapeString(u.CreatedAt.Format(adminDateFormat)),
		html.EscapeString(u.UpdatedAt.Format(adminDateFormat)),
		guid,
	)
}

func adminInvitesSection(invites []entity.AdminInviteView) string {
	rows := ""
	for _, inv := range invites {
		rows += adminInviteRow(inv)
	}
	return adminListSection(adminSectionHeader("Invites", "Create", "/admin/modals/invite"), fmt.Sprintf(`
<thead><tr><th>Guid</th><th>Team</th><th>Expired At</th><th>Actions</th></tr></thead>
<tbody>%s</tbody>`, rows))
}

func adminInviteRow(inv entity.AdminInviteView) string {
	guid := inv.Guid.String()
	return fmt.Sprintf(`<tr id="invite-%s">
	<td class="guid-cell">%s <button type="button" onclick="copyText('%s')">Copy</button></td>
	<td>%s</td>
	<td>%s</td>
	<td class="actions">
		<button type="button" onclick="copyInviteLink('%s')">Copy link</button>
		<button class="btn-danger" hx-delete="/admin/invites/%s" hx-target="closest tr" hx-swap="outerHTML settle:0.3s" hx-confirm="Delete invite?">Delete</button>
	</td>
</tr>`,
		guid, html.EscapeString(guid), guid,
		html.EscapeString(inv.TeamName),
		html.EscapeString(inv.ExpiredAt.Format(adminDateFormat)),
		guid, guid,
	)
}

func adminTeamsSection(teams []entity.Team) string {
	rows := ""
	for _, t := range teams {
		rows += adminTeamRow(t)
	}
	return adminListSection(adminSectionHeader("Teams", "Create", "/admin/modals/team"), fmt.Sprintf(`
<thead><tr><th>Name</th><th>Description</th><th>Actions</th></tr></thead>
<tbody>%s</tbody>`, rows))
}

func adminTeamRow(t entity.Team) string {
	return fmt.Sprintf(`<tr id="team-%s">
	<td>%s</td><td>%s</td>
	<td class="actions"><button class="btn-danger" hx-delete="/admin/teams/%s" hx-target="closest tr" hx-swap="outerHTML settle:0.3s" hx-confirm="Delete team?">Delete</button></td>
</tr>`,
		t.Guid, html.EscapeString(t.Name), html.EscapeString(strVal(t.Description)),
		t.Guid,
	)
}

func adminTeammatesSection(teammates []entity.AdminTeammateView) string {
	rows := ""
	for _, tm := range teammates {
		rows += adminTeammateRow(tm)
	}
	return adminListSection(adminSectionHeader("Teammates", "Invite", "/admin/modals/invite"), fmt.Sprintf(`
<thead><tr><th>Team</th><th>User</th><th>Leader</th><th>Actions</th></tr></thead>
<tbody>%s</tbody>`, rows))
}

func adminTeammateRow(tm entity.AdminTeammateView) string {
	return fmt.Sprintf(`<tr id="teammate-%s-%s">
	<td>%s</td><td>%s</td>
	<td>%s</td>
	<td class="actions">
		<button class="btn-danger" hx-delete="/admin/teammates/%s/%s" hx-target="closest tr" hx-swap="outerHTML settle:0.3s" hx-confirm="Delete teammate?">Delete</button>
	</td>
</tr>`,
		tm.UserGuid, tm.TeamGuid,
		html.EscapeString(tm.TeamName), html.EscapeString(tm.UserName),
		adminBoolCheckboxForm(
			fmt.Sprintf("/admin/teammates/%s/%s/leader", tm.UserGuid, tm.TeamGuid),
			"is_leader",
			tm.IsLeader,
		),
		tm.UserGuid, tm.TeamGuid,
	)
}

func adminWorkspacesSection(workspaces []entity.AdminWorkspaceView) string {
	rows := ""
	for _, w := range workspaces {
		rows += adminWorkspaceRow(w)
	}
	return adminListSection(adminSectionHeader("Workspaces", "", ""), fmt.Sprintf(`
<thead><tr><th>Name</th><th>Version</th><th>Team</th><th>Actions</th></tr></thead>
<tbody>%s</tbody>`, rows))
}

func adminWorkspaceRow(w entity.AdminWorkspaceView) string {
	return fmt.Sprintf(`<tr id="workspace-%s">
	<td>%s</td><td>%s</td><td>%s</td>
	<td class="actions"><button class="btn-danger" hx-delete="/admin/workspaces/%s" hx-target="closest tr" hx-swap="outerHTML settle:0.3s" hx-confirm="Delete workspace?">Delete</button></td>
</tr>`,
		w.Guid, html.EscapeString(w.Name), html.EscapeString(w.Version), html.EscapeString(w.TeamName),
		w.Guid,
	)
}

func adminWorkspaceRolesSection(roles []entity.AdminWorkspaceRoleView) string {
	rows := ""
	for _, r := range roles {
		rows += adminWorkspaceRoleRow(r)
	}
	return adminListSection(adminSectionHeader("Workspace Roles", "Create", "/admin/modals/workspace-role"), fmt.Sprintf(`
<thead><tr><th>Team</th><th>Workspace</th><th>User</th><th>Role</th><th>Actions</th></tr></thead>
<tbody>%s</tbody>`, rows))
}

func adminWorkspaceRoleRow(r entity.AdminWorkspaceRoleView) string {
	return fmt.Sprintf(`<tr id="wr-%s-%s">
	<td>%s</td><td>%s</td><td>%s</td><td>%s</td>
	<td class="actions"><button class="btn-danger" hx-delete="/admin/workspace_roles/%s/%s" hx-target="closest tr" hx-swap="outerHTML settle:0.3s" hx-confirm="Delete workspace role?">Delete</button></td>
</tr>`,
		r.WorkspaceGuid, r.TeammateGuid,
		html.EscapeString(r.TeamName), html.EscapeString(r.WorkspaceName), html.EscapeString(r.UserName), html.EscapeString(string(r.Role)),
		r.WorkspaceGuid, r.TeammateGuid,
	)
}

func adminInviteModal(teams []entity.Team, selectedTeamGUID string) string {
	options := teamSelectOptions(teams, selectedTeamGUID)
	return fmt.Sprintf(`<div class="modal-backdrop">
	<div class="modal-card">
		<h2>Create Invite</h2>
		<form hx-post="/admin/invites" hx-target="#content" hx-swap="innerHTML">
			<label for="team_guid">Team</label>
			<select id="team_guid" name="team_guid" required>%s</select>
			<div class="modal-actions">
				<button type="button" onclick="document.getElementById('modal-host').innerHTML=''">Cancel</button>
				<button type="submit" class="btn-primary">Create</button>
			</div>
		</form>
	</div>
</div>`, options)
}

func adminTeamModal() string {
	return `<div class="modal-backdrop">
	<div class="modal-card">
		<h2>Create Team</h2>
		<form hx-post="/admin/teams" hx-target="#content" hx-swap="innerHTML">
			<label for="name">Name</label>
			<input id="name" name="name" type="text" required maxlength="255">
			<label for="description">Description</label>
			<textarea id="description" name="description" rows="3"></textarea>
			<div class="modal-actions">
				<button type="button" onclick="document.getElementById('modal-host').innerHTML=''">Cancel</button>
				<button type="submit" class="btn-primary">Create</button>
			</div>
		</form>
	</div>
</div>`
}

func adminWorkspaceRoleModal(workspaces []entity.AdminWorkspaceView, users []entity.AdminUserView) string {
	wsOptions := `<option value="">Select workspace</option>`
	for _, w := range workspaces {
		wsOptions += fmt.Sprintf(`<option value="%s">%s (%s)</option>`,
			html.EscapeString(w.Guid.String()),
			html.EscapeString(w.Name),
			html.EscapeString(w.TeamName),
		)
	}
	userOptions := `<option value="">Select user</option>`
	for _, u := range users {
		label := u.Login
		if u.Name != nil && *u.Name != "" {
			label = *u.Name + " (" + u.Login + ")"
		}
		userOptions += fmt.Sprintf(`<option value="%s">%s</option>`,
			html.EscapeString(u.Guid.String()),
			html.EscapeString(label),
		)
	}
	return fmt.Sprintf(`<div class="modal-backdrop">
	<div class="modal-card">
		<h2>Create Workspace Role</h2>
		<form hx-post="/admin/workspace_roles" hx-target="#content" hx-swap="innerHTML">
			<label for="workspace_guid">Workspace</label>
			<select id="workspace_guid" name="workspace_guid" required>%s</select>
			<label for="teammate_guid">User</label>
			<select id="teammate_guid" name="teammate_guid" required>%s</select>
			<label for="role">Role</label>
			<select id="role" name="role" required>
				<option value="MAINTAINER">MAINTAINER</option>
				<option value="DEVELOPER">DEVELOPER</option>
				<option value="GUEST">GUEST</option>
			</select>
			<div class="modal-actions">
				<button type="button" onclick="document.getElementById('modal-host').innerHTML=''">Cancel</button>
				<button type="submit" class="btn-primary">Create</button>
			</div>
		</form>
	</div>
</div>`, wsOptions, userOptions)
}

func teamSelectOptions(teams []entity.Team, selected string) string {
	options := `<option value="">Select team</option>`
	for _, t := range teams {
		selectedAttr := ""
		if selected != "" && t.Guid.String() == selected {
			selectedAttr = " selected"
		}
		options += fmt.Sprintf(`<option value="%s"%s>%s</option>`,
			html.EscapeString(t.Guid.String()),
			selectedAttr,
			html.EscapeString(t.Name),
		)
	}
	return options
}

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
