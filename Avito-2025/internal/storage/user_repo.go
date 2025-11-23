package storage

import (
	"avito-2025/internal/domain"
	"context"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, team_id, is_active, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id`
	return r.db.QueryRowContext(ctx, query, user.Username, user.TeamID, user.IsActive).Scan(&user.ID)
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, team_id, is_active, created_at FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.TeamID, &user.IsActive, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // не найден
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetActiveMembers(ctx context.Context, teamID int) ([]*domain.User, error) {
	query := `SELECT id, username, team_id, is_active, created_at FROM users WHERE team_id = $1 AND is_active = TRUE`
	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Username, &user.TeamID, &user.IsActive, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET username=$1, is_active=$2 WHERE id=$3`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.IsActive, user.ID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
