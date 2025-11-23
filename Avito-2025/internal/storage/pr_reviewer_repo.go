package storage

import (
	"avito-2025/internal/domain"
	"context"
	"database/sql"
)

type PRReviewerRepository struct {
	db *sql.DB
}

func NewPRReviewerRepository(db *sql.DB) *PRReviewerRepository {
	return &PRReviewerRepository{db: db}
}

func (r *PRReviewerRepository) AssignReviewer(ctx context.Context, prID int, reviewerID int) error {
	query := `INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES ($1, $2)
              ON CONFLICT (pr_id, reviewer_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, prID, reviewerID)
	return err
}

func (r *PRReviewerRepository) ListByPR(ctx context.Context, prID int) ([]*domain.PRReviewer, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, pr_id, reviewer_id, assigned_at FROM pr_reviewers WHERE pr_id = $1`, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviewers []*domain.PRReviewer
	for rows.Next() {
		var rec domain.PRReviewer
		if err := rows.Scan(&rec.ID, &rec.PRID, &rec.ReviewerID, &rec.AssignedAt); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, &rec)
	}
	return reviewers, nil
}

func (r *PRReviewerRepository) RemoveReviewer(ctx context.Context, prID int, reviewerID int) error {
	query := `DELETE FROM pr_reviewers WHERE pr_id = $1 AND reviewer_id = $2`
	_, err := r.db.ExecContext(ctx, query, prID, reviewerID)
	return err
}

func (r *PRReviewerRepository) GetReviewersCount(ctx context.Context, prID int) (int, error) {
	query := `SELECT COUNT(*) FROM pr_reviewers WHERE pr_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, prID).Scan(&count)
	return count, err
}

func (r *PRReviewerRepository) ListByReviewer(ctx context.Context, reviewerID int) ([]*domain.PRReviewer, error) {
	query := `SELECT id, pr_id, reviewer_id, assigned_at FROM pr_reviewers WHERE reviewer_id = $1`
	rows, err := r.db.QueryContext(ctx, query, reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []*domain.PRReviewer
	for rows.Next() {
		var rec domain.PRReviewer
		if err := rows.Scan(&rec.ID, &rec.PRID, &rec.ReviewerID, &rec.AssignedAt); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, &rec)
	}
	return reviewers, nil
}
