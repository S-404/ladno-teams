package repository

import (
	"database/sql"
	"fmt"
	"ladno-teams/internal/entity"

	"github.com/google/uuid"
)

type IUserRepository interface {
	Create(data entity.User) (*entity.User, error)
	FindByGuid(guid uuid.UUID) (*entity.User, error)
	FindByLogin(login string) (*entity.User, error)
	Update(data entity.User) (*entity.User, error)
	Delete(guid uuid.UUID) error
	List(limit, offset int) ([]entity.User, error)
	Count() (int, error)
}

type UserRepository struct {
	*BaseRepository
	table string
}

func NewUserRepository(baseRepo *BaseRepository) *UserRepository {
	return &UserRepository{
		BaseRepository: baseRepo,
		table:          "users",
	}
}

func (r *UserRepository) Create(data entity.User) (*entity.User, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (login, password, is_admin, is_blocked)
		VALUES (:login, :password, :is_admin, :is_blocked)
		RETURNING guid, login, password, is_admin, is_blocked, created_at, updated_at, deleted_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user entity.User
	if rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserRepository) FindByGuid(guid uuid.UUID) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf(`SELECT * FROM %s WHERE guid = $1 AND deleted_at IS NULL`, r.table)
	if err := r.db.Get(&user, query, guid); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByLogin(login string) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf(`SELECT * FROM %s WHERE login = $1 AND deleted_at IS NULL`, r.table)
	if err := r.db.Get(&user, query, login); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(data entity.User) (*entity.User, error) {
	query := fmt.Sprintf(`
		UPDATE %s
		SET login = :login, is_admin = :is_admin, is_blocked = :is_blocked, updated_at = NOW()
		WHERE guid = :guid AND deleted_at IS NULL
		RETURNING guid, login, password, is_admin, is_blocked, created_at, updated_at, deleted_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user entity.User
	if rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
	} else {
		return nil, sql.ErrNoRows
	}

	return &user, nil
}

func (r *UserRepository) Delete(guid uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET deleted_at = NOW() WHERE guid = $1 AND deleted_at IS NULL`, r.table)
	_, err := r.db.Exec(query, guid)
	return err
}

func (r *UserRepository) List(limit, offset int) ([]entity.User, error) {
	var users []entity.User
	query := fmt.Sprintf(`SELECT * FROM %s WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`, r.table)
	if err := r.db.Select(&users, query, limit, offset); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Count() (int, error) {
	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE deleted_at IS NULL`, r.table)
	if err := r.db.Get(&count, query); err != nil {
		return 0, err
	}
	return count, nil
}
