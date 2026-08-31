package repository

import (
	"fmt"
	"ladno-teams/internal/entity"

	"github.com/google/uuid"
)

type ITeamRepository interface {
	Create(data entity.Team) (*entity.Team, error)
	FindByGuid(guid uuid.UUID) (*entity.Team, error)
	Update(data entity.Team) (*entity.Team, error)
	Delete(guid uuid.UUID) error
	ListByUser(userGuid uuid.UUID, limit, offset int) ([]entity.Team, error)
	ListAll(limit, offset int) ([]entity.Team, error)
	Count() (int, error)
}

type TeamRepository struct {
	*BaseRepository
	table string
}

func NewTeamRepository(baseRepo *BaseRepository) *TeamRepository {
	return &TeamRepository{
		BaseRepository: baseRepo,
		table:          "teams",
	}
}

func (r *TeamRepository) Create(data entity.Team) (*entity.Team, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (name, description)
		VALUES (:name, :description)
		RETURNING guid, name, description, created_at, updated_at, deleted_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var team entity.Team
	if rows.Next() {
		if err := rows.StructScan(&team); err != nil {
			return nil, err
		}
	}
	return &team, nil
}

func (r *TeamRepository) FindByGuid(guid uuid.UUID) (*entity.Team, error) {
	var team entity.Team
	query := fmt.Sprintf(`SELECT * FROM %s WHERE guid = $1 AND deleted_at IS NULL`, r.table)
	if err := r.db.Get(&team, query, guid); err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) Update(data entity.Team) (*entity.Team, error) {
	query := fmt.Sprintf(`
		UPDATE %s
		SET name = :name, description = :description, updated_at = NOW()
		WHERE guid = :guid AND deleted_at IS NULL
		RETURNING guid, name, description, created_at, updated_at, deleted_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var team entity.Team
	if rows.Next() {
		if err := rows.StructScan(&team); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("team not found")
	}
	return &team, nil
}

func (r *TeamRepository) Delete(guid uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET deleted_at = NOW() WHERE guid = $1 AND deleted_at IS NULL`, r.table)
	_, err := r.db.Exec(query, guid)
	return err
}

func (r *TeamRepository) ListByUser(userGuid uuid.UUID, limit, offset int) ([]entity.Team, error) {
	var teams []entity.Team
	query := fmt.Sprintf(`
		SELECT t.* FROM %s t
		JOIN teammates tm ON t.guid = tm.team_guid
		WHERE t.deleted_at IS NULL AND tm.user_guid = $1
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3`, r.table)
	if err := r.db.Select(&teams, query, userGuid, limit, offset); err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *TeamRepository) ListAll(limit, offset int) ([]entity.Team, error) {
	var teams []entity.Team
	query := fmt.Sprintf(`SELECT * FROM %s WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`, r.table)
	if err := r.db.Select(&teams, query, limit, offset); err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *TeamRepository) Count() (int, error) {
	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE deleted_at IS NULL`, r.table)
	if err := r.db.Get(&count, query); err != nil {
		return 0, err
	}
	return count, nil
}
