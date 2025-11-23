package storage

import (
	"avito-2025/internal/domain"
	"context"
	"database/sql"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(ctx context.Context, team *domain.Team) error {
	query := `INSERT INTO teams (name) VALUES ($1) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query, team.Name).Scan(&team.ID, &team.CreatedAt)
}

func (r *TeamRepository) GetByID(ctx context.Context, id int) (*domain.Team, error) {
	var team domain.Team
	query := `SELECT id, name, created_at FROM teams WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&team.ID, &team.Name, &team.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &team, err
}

func (r *TeamRepository) List(ctx context.Context) ([]*domain.Team, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, created_at FROM teams`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teams []*domain.Team
	for rows.Next() {
		var team domain.Team
		if err := rows.Scan(&team.ID, &team.Name, &team.CreatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, &team)
	}
	return teams, nil
}

func (r *TeamRepository) Update(ctx context.Context, team *domain.Team) error {
	query := `UPDATE teams SET name = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, team.Name, team.ID)
	return err
}

func (r *TeamRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM teams WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
