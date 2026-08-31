package repository

import (
	"fmt"
	"ladno-teams/internal/entity"

	"github.com/google/uuid"
)

type IWorkspaceRepository interface {
	Create(data entity.Workspace) (*entity.Workspace, error)
	FindByGuid(guid uuid.UUID) (*entity.Workspace, error)
	Update(data entity.Workspace) (*entity.Workspace, error)
	Delete(guid uuid.UUID) error
	ListByTeam(teamGuid uuid.UUID, limit, offset int) ([]entity.Workspace, error)
	ListAll(limit, offset int) ([]entity.Workspace, error)
	Count() (int, error)
}

type WorkspaceRepository struct {
	*BaseRepository
	table string
}

func NewWorkspaceRepository(baseRepo *BaseRepository) *WorkspaceRepository {
	return &WorkspaceRepository{
		BaseRepository: baseRepo,
		table:          "workspaces",
	}
}

func (r *WorkspaceRepository) Create(data entity.Workspace) (*entity.Workspace, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (team_guid, name, version, data, config, envs)
		VALUES (:team_guid, :name, :version, :data, :config, :envs)
		RETURNING guid, team_guid, name, version, data, config, envs, created_at, updated_at, deleted_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspace entity.Workspace
	if rows.Next() {
		if err := rows.StructScan(&workspace); err != nil {
			return nil, err
		}
	}
	return &workspace, nil
}

func (r *WorkspaceRepository) FindByGuid(guid uuid.UUID) (*entity.Workspace, error) {
	var workspace entity.Workspace
	query := fmt.Sprintf(`SELECT * FROM %s WHERE guid = $1 AND deleted_at IS NULL`, r.table)
	if err := r.db.Get(&workspace, query, guid); err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (r *WorkspaceRepository) Update(data entity.Workspace) (*entity.Workspace, error) {
	query := fmt.Sprintf(`
		UPDATE %s
		SET name = :name, version = :version, data = :data, config = :config, envs = :envs, updated_at = NOW()
		WHERE guid = :guid AND deleted_at IS NULL
		RETURNING guid, team_guid, name, version, data, config, envs, created_at, updated_at, deleted_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspace entity.Workspace
	if rows.Next() {
		if err := rows.StructScan(&workspace); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("workspace not found")
	}
	return &workspace, nil
}

func (r *WorkspaceRepository) Delete(guid uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET deleted_at = NOW() WHERE guid = $1 AND deleted_at IS NULL`, r.table)
	_, err := r.db.Exec(query, guid)
	return err
}

func (r *WorkspaceRepository) ListByTeam(teamGuid uuid.UUID, limit, offset int) ([]entity.Workspace, error) {
	var workspaces []entity.Workspace
	query := fmt.Sprintf(`
		SELECT * FROM %s
		WHERE team_guid = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`, r.table)
	if err := r.db.Select(&workspaces, query, teamGuid, limit, offset); err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (r *WorkspaceRepository) ListAll(limit, offset int) ([]entity.Workspace, error) {
	var workspaces []entity.Workspace
	query := fmt.Sprintf(`SELECT * FROM %s WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`, r.table)
	if err := r.db.Select(&workspaces, query, limit, offset); err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (r *WorkspaceRepository) Count() (int, error) {
	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE deleted_at IS NULL`, r.table)
	if err := r.db.Get(&count, query); err != nil {
		return 0, err
	}
	return count, nil
}
