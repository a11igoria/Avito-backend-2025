package storage

import (
	"context"
	"database/sql"
	"strconv"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create — создать пользователя
func (r *UserRepository) Create(ctx context.Context, username string, teamName string, isActive bool) (string, error) {
	var id int
	query := `INSERT INTO users (username, team_name, is_active, created_at) 
	          VALUES ($1, $2, $3, NOW()) 
	          RETURNING id`

	err := r.db.QueryRowContext(ctx, query, username, teamName, isActive).Scan(&id)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(id), nil
}

// GetByID — получить пользователя по string ID
func (r *UserRepository) GetByID(ctx context.Context, userID string) (map[string]interface{}, error) {
	var id int
	var username, teamName string
	var isActive bool
	var createdAt interface{}

	query := `SELECT id, username, team_name, is_active, created_at 
	          FROM users WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, userID).
		Scan(&id, &username, &teamName, &isActive, &createdAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ID":       strconv.Itoa(id),
		"Username": username,
		"TeamName": teamName,
		"IsActive": isActive,
	}, nil
}

// GetActiveMembers — получить активных членов команды по имени команды
func (r *UserRepository) GetActiveMembers(ctx context.Context, teamName string) ([]map[string]interface{}, error) {
	query := `SELECT id, username, team_name, is_active, created_at 
	          FROM users WHERE team_name = $1 AND is_active = TRUE`

	rows, err := r.db.QueryContext(ctx, query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var username, teamName string
		var isActive bool
		var createdAt interface{}

		if err := rows.Scan(&id, &username, &teamName, &isActive, &createdAt); err != nil {
			return nil, err
		}

		users = append(users, map[string]interface{}{
			"ID":       strconv.Itoa(id),
			"Username": username,
			"TeamName": teamName,
			"IsActive": isActive,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Update — обновить пользователя
func (r *UserRepository) Update(ctx context.Context, userID string, username string, isActive bool) error {
	query := `UPDATE users SET username=$1, is_active=$2, updated_at=NOW() WHERE id=$3`
	_, err := r.db.ExecContext(ctx, query, username, isActive, userID)
	return err
}

// Delete — удалить пользователя
func (r *UserRepository) Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
