package storage

import (
	"context"
	"database/sql"
	"strconv"
)

type PRReviewerRepository struct {
	db *sql.DB
}

func NewPRReviewerRepository(db *sql.DB) *PRReviewerRepository {
	return &PRReviewerRepository{db: db}
}

// AssignReviewer — назначить ревьювера на PR
func (r *PRReviewerRepository) AssignReviewer(ctx context.Context, prID string, reviewerID string) error {
	query := `INSERT INTO pr_reviewers (pr_id, reviewer_id, assigned_at) 
	          VALUES ($1, $2, NOW())
	          ON CONFLICT (pr_id, reviewer_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, prID, reviewerID)
	return err
}

// GetByPR — получить всех ревьюверов PR по string ID
func (r *PRReviewerRepository) GetByPR(ctx context.Context, prID string) ([]string, error) {
	query := `SELECT reviewer_id FROM pr_reviewers WHERE pr_id = $1 ORDER BY assigned_at`

	rows, err := r.db.QueryContext(ctx, query, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewerIDs []string
	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, err
		}
		reviewerIDs = append(reviewerIDs, reviewerID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviewerIDs, nil
}

// RemoveReviewer — удалить ревьювера с PR
func (r *PRReviewerRepository) RemoveReviewer(ctx context.Context, prID string, reviewerID string) error {
	query := `DELETE FROM pr_reviewers WHERE pr_id = $1 AND reviewer_id = $2`
	_, err := r.db.ExecContext(ctx, query, prID, reviewerID)
	return err
}

// GetReviewersCount — получить количество ревьюверов для PR
func (r *PRReviewerRepository) GetReviewersCount(ctx context.Context, prID string) (int, error) {
	query := `SELECT COUNT(*) FROM pr_reviewers WHERE pr_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, prID).Scan(&count)
	return count, err
}

// GetPRsByReviewer — получить все PR где юзер ревьювер
func (r *PRReviewerRepository) GetPRsByReviewer(ctx context.Context, reviewerID string) ([]string, error) {
	query := `SELECT DISTINCT pr_id FROM pr_reviewers WHERE reviewer_id = $1 ORDER BY pr_id`

	rows, err := r.db.QueryContext(ctx, query, reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prIDs []string
	for rows.Next() {
		var prID string
		if err := rows.Scan(&prID); err != nil {
			return nil, err
		}
		prIDs = append(prIDs, prID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prIDs, nil
}

// List — получить все связи ревьювер-PR
func (r *PRReviewerRepository) List(ctx context.Context) ([]map[string]interface{}, error) {
	query := `SELECT id, pr_id, reviewer_id, assigned_at FROM pr_reviewers ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []map[string]interface{}
	for rows.Next() {
		var id int
		var prID, reviewerID string
		var assignedAt interface{}

		if err := rows.Scan(&id, &prID, &reviewerID, &assignedAt); err != nil {
			return nil, err
		}

		reviewers = append(reviewers, map[string]interface{}{
			"ID":         strconv.Itoa(id),
			"PRID":       prID,
			"ReviewerID": reviewerID,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviewers, nil
}
