package repository

import (
	"database/sql"
	"library-project/internal/models"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(rt *models.RefreshToken) error {
	rt.ID = uuid.New().String()

	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`

	return r.db.QueryRow(query, rt.ID, rt.UserID, rt.Token, rt.ExpiresAt).
		Scan(&rt.CreatedAt)
}

func (r *RefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	rt := &models.RefreshToken{}

	query := `SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = $1`

	err := r.db.QueryRow(query, token).Scan(
		&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return rt, err
}

func (r *RefreshTokenRepository) DeleteByToken(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.Exec(query, token)
	return err
}

func (r *RefreshTokenRepository) DeleteByUserID(userID string) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *RefreshTokenRepository) DeleteExpired() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`
	_, err := r.db.Exec(query, time.Now())
	return err
}
