package repository

import (
	"fmt"
	"ladno-teams/internal/entity"

	"github.com/google/uuid"
)

type IWorkspaceRoleRepository interface {
	Create(data entity.WorkspaceRole) (*entity.WorkspaceRole, error)
	FindByWorkspaceAndTeammate(workspaceGuid, teammateGuid uuid.UUID) (*entity.WorkspaceRole, error)
	FindByWorkspace(workspaceGuid uuid.UUID) ([]entity.WorkspaceRole, error)
	Update(data entity.WorkspaceRole) (*entity.WorkspaceRole, error)
	Delete(workspaceGuid, teammateGuid uuid.UUID) error
}

type WorkspaceRoleRepository struct {
	*BaseRepository
	table string
}

func NewWorkspaceRoleRepository(baseRepo *BaseRepository) *WorkspaceRoleRepository {
	return &WorkspaceRoleRepository{
		BaseRepository: baseRepo,
		table:          "workspace_roles",
	}
}

func (r *WorkspaceRoleRepository) Create(data entity.WorkspaceRole) (*entity.WorkspaceRole, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (workspace_guid, teammate_guid, role)
		VALUES (:workspace_guid, :teammate_guid, :role)
		RETURNING workspace_guid, teammate_guid, role, created_at, updated_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var role entity.WorkspaceRole
	if rows.Next() {
		if err := rows.StructScan(&role); err != nil {
			return nil, err
		}
	}
	return &role, nil
}

func (r *WorkspaceRoleRepository) FindByWorkspaceAndTeammate(workspaceGuid, teammateGuid uuid.UUID) (*entity.WorkspaceRole, error) {
	var role entity.WorkspaceRole
	query := fmt.Sprintf(`SELECT * FROM %s WHERE workspace_guid = $1 AND teammate_guid = $2`, r.table)
	if err := r.db.Get(&role, query, workspaceGuid, teammateGuid); err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *WorkspaceRoleRepository) FindByWorkspace(workspaceGuid uuid.UUID) ([]entity.WorkspaceRole, error) {
	var roles []entity.WorkspaceRole
	query := fmt.Sprintf(`SELECT * FROM %s WHERE workspace_guid = $1 ORDER BY created_at`, r.table)
	if err := r.db.Select(&roles, query, workspaceGuid); err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *WorkspaceRoleRepository) Update(data entity.WorkspaceRole) (*entity.WorkspaceRole, error) {
	query := fmt.Sprintf(`
		UPDATE %s
		SET role = :role, updated_at = NOW()
		WHERE workspace_guid = :workspace_guid AND teammate_guid = :teammate_guid
		RETURNING workspace_guid, teammate_guid, role, created_at, updated_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var role entity.WorkspaceRole
	if rows.Next() {
		if err := rows.StructScan(&role); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("workspace role not found")
	}
	return &role, nil
}

func (r *WorkspaceRoleRepository) Delete(workspaceGuid, teammateGuid uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE workspace_guid = $1 AND teammate_guid = $2`, r.table)
	_, err := r.db.Exec(query, workspaceGuid, teammateGuid)
	return err
}
