package repository

import (
	"fmt"
	"ladno-teams/internal/entity"

	"github.com/google/uuid"
)

type IProfileRepository interface {
	Create(data entity.Profile) (*entity.Profile, error)
	FindByUserGuid(userGuid uuid.UUID) (*entity.Profile, error)
	Update(data entity.Profile) (*entity.Profile, error)
	DeleteByUserGuid(userGuid uuid.UUID) error
}

type ProfileRepository struct {
	*BaseRepository
	table string
}

func NewProfileRepository(baseRepo *BaseRepository) *ProfileRepository {
	return &ProfileRepository{
		BaseRepository: baseRepo,
		table:          "profiles",
	}
}

func (r *ProfileRepository) Create(data entity.Profile) (*entity.Profile, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_guid, name, about)
		VALUES (:user_guid, :name, :about)
		RETURNING user_guid, name, about, created_at, updated_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profile entity.Profile
	if rows.Next() {
		if err := rows.StructScan(&profile); err != nil {
			return nil, err
		}
	}

	return &profile, nil
}

func (r *ProfileRepository) FindByUserGuid(userGuid uuid.UUID) (*entity.Profile, error) {
	var profile entity.Profile
	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_guid = $1`, r.table)
	if err := r.db.Get(&profile, query, userGuid); err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) Update(data entity.Profile) (*entity.Profile, error) {
	query := fmt.Sprintf(`
		UPDATE %s
		SET name = :name, about = :about, updated_at = NOW()
		WHERE user_guid = :user_guid
		RETURNING user_guid, name, about, created_at, updated_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profile entity.Profile
	if rows.Next() {
		if err := rows.StructScan(&profile); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("profile not found")
	}

	return &profile, nil
}

func (r *ProfileRepository) DeleteByUserGuid(userGuid uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_guid = $1`, r.table)
	_, err := r.db.Exec(query, userGuid)
	return err
}
