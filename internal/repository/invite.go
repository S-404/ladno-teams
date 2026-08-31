package repository

import (
	"fmt"
	"ladno-teams/internal/entity"
	"time"

	"github.com/google/uuid"
)

type IInviteRepository interface {
	Create(data entity.Invite) (*entity.Invite, error)
	FindByGuid(guid uuid.UUID) (*entity.Invite, error)
	Delete(guid uuid.UUID) error
	DeleteExpired() error
}

type InviteRepository struct {
	*BaseRepository
	table string
}

func NewInviteRepository(baseRepo *BaseRepository) *InviteRepository {
	return &InviteRepository{
		BaseRepository: baseRepo,
		table:          "invites",
	}
}

func (r *InviteRepository) Create(data entity.Invite) (*entity.Invite, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (team_guid, expired_at)
		VALUES (:team_guid, :expired_at)
		RETURNING guid, team_guid, created_at, expired_at`, r.table)

	rows, err := r.db.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invite entity.Invite
	if rows.Next() {
		if err := rows.StructScan(&invite); err != nil {
			return nil, err
		}
	}
	return &invite, nil
}

func (r *InviteRepository) FindByGuid(guid uuid.UUID) (*entity.Invite, error) {
	var invite entity.Invite
	query := fmt.Sprintf(`SELECT * FROM %s WHERE guid = $1 AND expired_at > $2`, r.table)
	if err := r.db.Get(&invite, query, guid, time.Now().UTC()); err != nil {
		return nil, err
	}
	return &invite, nil
}

func (r *InviteRepository) Delete(guid uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE guid = $1`, r.table)
	_, err := r.db.Exec(query, guid)
	return err
}

func (r *InviteRepository) DeleteExpired() error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE expired_at <= $1`, r.table)
	_, err := r.db.Exec(query, time.Now().UTC())
	return err
}
