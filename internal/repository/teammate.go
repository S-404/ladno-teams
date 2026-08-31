package repository

import (
	"fmt"
	"ladno-teams/internal/entity"

	"github.com/google/uuid"
)

type ITeammateRepository interface {
	Create(data entity.Teammate) (*entity.Teammate, error)
	FindByUserAndTeam(userGuid, teamGuid uuid.UUID) (*entity.Teammate, error)
	FindByTeam(teamGuid uuid.UUID) ([]entity.Teammate, error)
	Update(data entity.Teammate) (*entity.Teammate, error)
	Delete(userGuid, teamGuid uuid.UUID) error
	IsLeader(userGuid, teamGuid uuid.UUID) (bool, error)
	IsMember(userGuid, teamGuid uuid.UUID) (bool, error)
	HasLeader(teamGuid uuid.UUID) (bool, error)
}

type TeammateRepository struct {
	*BaseRepository
	table string
}

func NewTeammateRepository(baseRepo *BaseRepository) *TeammateRepository {
	return &TeammateRepository{
		BaseRepository: baseRepo,
		table:          "teammates",
	}
}

func (r *TeammateRepository) Create(data entity.Teammate) (*entity.Teammate, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_guid, team_guid, is_leader)
		VALUES (:user_guid, :team_guid, :is_leader)
		RETURNING user_guid, team_guid, is_leader, created_at, updated_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teammate entity.Teammate
	if rows.Next() {
		if err := rows.StructScan(&teammate); err != nil {
			return nil, err
		}
	}
	return &teammate, nil
}

func (r *TeammateRepository) FindByUserAndTeam(userGuid, teamGuid uuid.UUID) (*entity.Teammate, error) {
	var teammate entity.Teammate
	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_guid = $1 AND team_guid = $2`, r.table)
	if err := r.db.Get(&teammate, query, userGuid, teamGuid); err != nil {
		return nil, err
	}
	return &teammate, nil
}

func (r *TeammateRepository) FindByTeam(teamGuid uuid.UUID) ([]entity.Teammate, error) {
	var teammates []entity.Teammate
	query := fmt.Sprintf(`SELECT * FROM %s WHERE team_guid = $1 ORDER BY created_at`, r.table)
	if err := r.db.Select(&teammates, query, teamGuid); err != nil {
		return nil, err
	}
	return teammates, nil
}

func (r *TeammateRepository) Update(data entity.Teammate) (*entity.Teammate, error) {
	query := fmt.Sprintf(`
		UPDATE %s
		SET is_leader = :is_leader, updated_at = NOW()
		WHERE user_guid = :user_guid AND team_guid = :team_guid
		RETURNING user_guid, team_guid, is_leader, created_at, updated_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teammate entity.Teammate
	if rows.Next() {
		if err := rows.StructScan(&teammate); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("teammate not found")
	}
	return &teammate, nil
}

func (r *TeammateRepository) Delete(userGuid, teamGuid uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_guid = $1 AND team_guid = $2`, r.table)
	_, err := r.db.Exec(query, userGuid, teamGuid)
	return err
}

func (r *TeammateRepository) IsLeader(userGuid, teamGuid uuid.UUID) (bool, error) {
	var exists bool
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE user_guid = $1 AND team_guid = $2 AND is_leader = TRUE)`, r.table)
	if err := r.db.Get(&exists, query, userGuid, teamGuid); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *TeammateRepository) IsMember(userGuid, teamGuid uuid.UUID) (bool, error) {
	var exists bool
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE user_guid = $1 AND team_guid = $2)`, r.table)
	if err := r.db.Get(&exists, query, userGuid, teamGuid); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *TeammateRepository) HasLeader(teamGuid uuid.UUID) (bool, error) {
	var exists bool
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE team_guid = $1 AND is_leader = TRUE)`, r.table)
	if err := r.db.Get(&exists, query, teamGuid); err != nil {
		return false, err
	}
	return exists, nil
}
