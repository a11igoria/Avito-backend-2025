package storage

import (
	"avito-2025/internal/domain"
	"context"
	"database/sql"
)

type PRRepository struct {
	db *sql.DB
}

func NewPRRepository(db *sql.DB) *PRRepository {
	return &PRRepository{db: db}
}

func (r *PRRepository) Create(ctx context.Context, pr *domain.PullRequest) error {
	query := `INSERT INTO pull_requests (title, author_id, status) 
              VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, pr.Title, pr.AuthorID, pr.Status).
		Scan(&pr.ID, &pr.CreatedAt, &pr.UpdatedAt)
}

func (r *PRRepository) GetByID(ctx context.Context, id int) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	query := `SELECT id, title, author_id, status, created_at, updated_at FROM pull_requests WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &pr, err
}

func (r *PRRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE pull_requests SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *PRRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM pull_requests WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PRRepository) List(ctx context.Context) ([]*domain.PullRequest, error) {
	query := `SELECT id, title, author_id, status, created_at, updated_at FROM pull_requests`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*domain.PullRequest
	for rows.Next() {
		var pr domain.PullRequest
		if err := rows.Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.UpdatedAt); err != nil {
			return nil, err
		}
		prs = append(prs, &pr)
	}
	return prs, nil
}

func (r *PRRepository) GetByAuthor(ctx context.Context, authorID int) ([]*domain.PullRequest, error) {
	query := `SELECT id, title, author_id, status, created_at, updated_at FROM pull_requests WHERE author_id = $1`
	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*domain.PullRequest
	for rows.Next() {
		var pr domain.PullRequest
		if err := rows.Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.UpdatedAt); err != nil {
			return nil, err
		}
		prs = append(prs, &pr)
	}
	return prs, nil
}

func (r *PRRepository) GetByStatus(ctx context.Context, status string) ([]*domain.PullRequest, error) {
	query := `SELECT id, title, author_id, status, created_at, updated_at FROM pull_requests WHERE status = $1`
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*domain.PullRequest
	for rows.Next() {
		var pr domain.PullRequest
		if err := rows.Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.UpdatedAt); err != nil {
			return nil, err
		}
		prs = append(prs, &pr)
	}
	return prs, nil
}
