package repository

import "ladno-teams/internal/entity"

type IAdminRepository interface {
	ListUsersWithProfiles(limit, offset int) ([]entity.AdminUserView, error)
	ListInvitesWithTeam(limit, offset int) ([]entity.AdminInviteView, error)
	ListTeammatesWithDetails(limit, offset int) ([]entity.AdminTeammateView, error)
	ListWorkspacesWithTeam(limit, offset int) ([]entity.AdminWorkspaceView, error)
	ListWorkspaceRolesWithDetails(limit, offset int) ([]entity.AdminWorkspaceRoleView, error)
}

type AdminRepository struct {
	*BaseRepository
}

func NewAdminRepository(baseRepo *BaseRepository) *AdminRepository {
	return &AdminRepository{BaseRepository: baseRepo}
}

func (r *AdminRepository) ListUsersWithProfiles(limit, offset int) ([]entity.AdminUserView, error) {
	var rows []entity.AdminUserView
	query := `
		SELECT u.guid, u.login, u.is_admin, u.is_blocked, u.created_at, u.updated_at,
		       p.name, p.about
		FROM users u
		LEFT JOIN profiles p ON p.user_guid = u.guid
		WHERE u.deleted_at IS NULL
		ORDER BY u.created_at DESC
		LIMIT $1 OFFSET $2`
	if err := r.db.Select(&rows, query, limit, offset); err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *AdminRepository) ListInvitesWithTeam(limit, offset int) ([]entity.AdminInviteView, error) {
	var rows []entity.AdminInviteView
	query := `
		SELECT i.guid, i.team_guid, t.name AS team_name, i.expired_at
		FROM invites i
		JOIN teams t ON t.guid = i.team_guid
		WHERE t.deleted_at IS NULL
		ORDER BY i.created_at DESC
		LIMIT $1 OFFSET $2`
	if err := r.db.Select(&rows, query, limit, offset); err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *AdminRepository) ListTeammatesWithDetails(limit, offset int) ([]entity.AdminTeammateView, error) {
	var rows []entity.AdminTeammateView
	query := `
		SELECT tm.user_guid, tm.team_guid, tm.is_leader,
		       t.name AS team_name,
		       COALESCE(p.name, u.login) AS user_name
		FROM teammates tm
		JOIN teams t ON t.guid = tm.team_guid
		JOIN users u ON u.guid = tm.user_guid
		LEFT JOIN profiles p ON p.user_guid = u.guid
		WHERE t.deleted_at IS NULL AND u.deleted_at IS NULL
		ORDER BY tm.created_at DESC
		LIMIT $1 OFFSET $2`
	if err := r.db.Select(&rows, query, limit, offset); err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *AdminRepository) ListWorkspacesWithTeam(limit, offset int) ([]entity.AdminWorkspaceView, error) {
	var rows []entity.AdminWorkspaceView
	query := `
		SELECT w.guid, w.name, w.version, t.name AS team_name
		FROM workspaces w
		JOIN teams t ON t.guid = w.team_guid
		WHERE w.deleted_at IS NULL AND t.deleted_at IS NULL
		ORDER BY w.created_at DESC
		LIMIT $1 OFFSET $2`
	if err := r.db.Select(&rows, query, limit, offset); err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *AdminRepository) ListWorkspaceRolesWithDetails(limit, offset int) ([]entity.AdminWorkspaceRoleView, error) {
	var rows []entity.AdminWorkspaceRoleView
	query := `
		SELECT wr.workspace_guid, wr.teammate_guid, wr.role,
		       t.name AS team_name,
		       w.name AS workspace_name,
		       COALESCE(p.name, u.login) AS user_name
		FROM workspace_roles wr
		JOIN workspaces w ON w.guid = wr.workspace_guid
		JOIN teams t ON t.guid = w.team_guid
		JOIN users u ON u.guid = wr.teammate_guid
		LEFT JOIN profiles p ON p.user_guid = u.guid
		WHERE w.deleted_at IS NULL AND t.deleted_at IS NULL AND u.deleted_at IS NULL
		ORDER BY wr.created_at DESC
		LIMIT $1 OFFSET $2`
	if err := r.db.Select(&rows, query, limit, offset); err != nil {
		return nil, err
	}
	return rows, nil
}
