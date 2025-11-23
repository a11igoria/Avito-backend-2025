package storage

import (
	"context"
	"database/sql"
	"strconv"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

// Create — создать команду
func (r *TeamRepository) Create(ctx context.Context, teamName string) error {
	query := `INSERT INTO teams (name, created_at) VALUES ($1, NOW())`
	_, err := r.db.ExecContext(ctx, query, teamName)
	return err
}

// GetByName — получить команду по имени (string)
func (r *TeamRepository) GetByName(ctx context.Context, teamName string) (map[string]interface{}, error) {
	var id int
	var name string
	var createdAt interface{}

	query := `SELECT id, name, created_at FROM teams WHERE name = $1`
	err := r.db.QueryRowContext(ctx, query, teamName).
		Scan(&id, &name, &createdAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ID":   strconv.Itoa(id),
		"Name": name,
	}, nil
}

// List — получить все команды
func (r *TeamRepository) List(ctx context.Context) ([]map[string]interface{}, error) {
	query := `SELECT id, name, created_at FROM teams ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var createdAt interface{}

		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			return nil, err
		}

		teams = append(teams, map[string]interface{}{
			"ID":   strconv.Itoa(id),
			"Name": name,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

// Update — обновить команду по имени
func (r *TeamRepository) Update(ctx context.Context, oldTeamName string, newTeamName string) error {
	query := `UPDATE teams SET name = $1 WHERE name = $2`
	_, err := r.db.ExecContext(ctx, query, newTeamName, oldTeamName)
	return err
}

// Delete — удалить команду по имени
func (r *TeamRepository) Delete(ctx context.Context, teamName string) error {
	query := `DELETE FROM teams WHERE name = $1`
	_, err := r.db.ExecContext(ctx, query, teamName)
	return err
}
