package storage

import (
	"context"
	"database/sql"
	"strconv"
)

type PRRepository struct {
	db *sql.DB
}

func NewPRRepository(db *sql.DB) *PRRepository {
	return &PRRepository{db: db}
}

// Create — создать PR
func (r *PRRepository) Create(ctx context.Context, prName string, authorID string, status string) (string, error) {
	var id int
	query := `INSERT INTO pull_requests (name, author_id, status, created_at) 
	          VALUES ($1, $2, $3, NOW()) 
	          RETURNING id`

	err := r.db.QueryRowContext(ctx, query, prName, authorID, status).Scan(&id)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(id), nil
}

// GetByID — получить PR по string ID (конвертирует строку в число для SQL!)
func (r *PRRepository) GetByID(ctx context.Context, prID string) (map[string]interface{}, error) {
	// ✅ ВАЖНО: Конвертируем string ID в int для SQL запроса
	idInt, err := strconv.Atoi(prID)
	if err != nil {
		return nil, err
	}

	var id int
	var name, authorID, status string
	var createdAt, updatedAt interface{}

	query := `SELECT id, name, author_id, status, created_at, updated_at 
	          FROM pull_requests WHERE id = $1`

	err = r.db.QueryRowContext(ctx, query, idInt).
		Scan(&id, &name, &authorID, &status, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ID":       strconv.Itoa(id),
		"Name":     name,
		"AuthorID": authorID,
		"Status":   status,
	}, nil
}

// UpdateStatus — обновить статус PR
func (r *PRRepository) UpdateStatus(ctx context.Context, prID string, status string) error {
	// ✅ ВАЖНО: Конвертируем string ID в int для SQL запроса
	idInt, err := strconv.Atoi(prID)
	if err != nil {
		return err
	}

	query := `UPDATE pull_requests SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err = r.db.ExecContext(ctx, query, status, idInt)
	return err
}

// Delete — удалить PR
func (r *PRRepository) Delete(ctx context.Context, prID string) error {
	idInt, err := strconv.Atoi(prID)
	if err != nil {
		return err
	}

	query := `DELETE FROM pull_requests WHERE id = $1`
	_, err = r.db.ExecContext(ctx, query, idInt)
	return err
}

// List — получить все PR
func (r *PRRepository) List(ctx context.Context) ([]map[string]interface{}, error) {
	query := `SELECT id, name, author_id, status, created_at, updated_at 
	          FROM pull_requests ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []map[string]interface{}
	for rows.Next() {
		var id int
		var name, authorID, status string
		var createdAt, updatedAt interface{}

		if err := rows.Scan(&id, &name, &authorID, &status, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		prs = append(prs, map[string]interface{}{
			"ID":       strconv.Itoa(id),
			"Name":     name,
			"AuthorID": authorID,
			"Status":   status,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}

// GetByAuthorID — получить все PR автора
func (r *PRRepository) GetByAuthorID(ctx context.Context, authorID string) ([]map[string]interface{}, error) {
	query := `SELECT id, name, author_id, status, created_at, updated_at 
	          FROM pull_requests WHERE author_id = $1 ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []map[string]interface{}
	for rows.Next() {
		var id int
		var name, authorID, status string
		var createdAt, updatedAt interface{}

		if err := rows.Scan(&id, &name, &authorID, &status, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		prs = append(prs, map[string]interface{}{
			"ID":       strconv.Itoa(id),
			"Name":     name,
			"AuthorID": authorID,
			"Status":   status,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}

// GetByStatus — получить все PR с определенным статусом
func (r *PRRepository) GetByStatus(ctx context.Context, status string) ([]map[string]interface{}, error) {
	query := `SELECT id, name, author_id, status, created_at, updated_at 
	          FROM pull_requests WHERE status = $1 ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []map[string]interface{}
	for rows.Next() {
		var id int
		var name, authorID, status string
		var createdAt, updatedAt interface{}

		if err := rows.Scan(&id, &name, &authorID, &status, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		prs = append(prs, map[string]interface{}{
			"ID":       strconv.Itoa(id),
			"Name":     name,
			"AuthorID": authorID,
			"Status":   status,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}
