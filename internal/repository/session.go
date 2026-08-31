package repository

import (
	"fmt"
	"ladno-teams/internal/entity"
	"time"

	"github.com/google/uuid"
)

type ISessionRepository interface {
	Create(data entity.Session) (*entity.Session, error)
	FindByToken(token string) (*entity.Session, error)
	FindByUserID(userID uuid.UUID) (*entity.Session, error)
	DeleteByToken(token string) error
	DeleteByUserID(userID uuid.UUID) error
	DeleteExpired() error
}

type SessionRepository struct {
	*BaseRepository
	table string
}

func NewSessionRepository(baseRepo *BaseRepository) *SessionRepository {
	return &SessionRepository{
		BaseRepository: baseRepo,
		table:          "sessions",
	}
}

func (r *SessionRepository) Create(data entity.Session) (*entity.Session, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, token, expired_at)
		VALUES (:user_id, :token, :expired_at)
		RETURNING guid, user_id, token, expired_at, created_at, updated_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var session entity.Session
	if rows.Next() {
		if err := rows.StructScan(&session); err != nil {
			return nil, err
		}
	}

	return &session, nil
}

func (r *SessionRepository) FindByToken(token string) (*entity.Session, error) {
	var session entity.Session
	query := fmt.Sprintf(`SELECT * FROM %s WHERE token = $1 AND expired_at > $2`, r.table)
	if err := r.db.Get(&session, query, token, time.Now().UTC()); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) FindByUserID(userID uuid.UUID) (*entity.Session, error) {
	var session entity.Session
	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1 AND expired_at > $2`, r.table)
	if err := r.db.Get(&session, query, userID, time.Now().UTC()); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) DeleteByToken(token string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE token = $1`, r.table)
	_, err := r.db.Exec(query, token)
	return err
}

func (r *SessionRepository) DeleteByUserID(userID uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = $1`, r.table)
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *SessionRepository) DeleteExpired() error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE expired_at <= $1`, r.table)
	_, err := r.db.Exec(query, time.Now().UTC())
	return err
}
