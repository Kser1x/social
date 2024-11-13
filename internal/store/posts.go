package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type PostModel struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserId    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *PostModel) error {
	query := `INSERT INTO posts (content, title, user_id, tags)
				VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
`
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserId,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err

	}
	return nil
}
func (s *PostsStore) GetByID(ctx context.Context, id int64) (*PostModel, error) {
	query := `
SELECT id, user_id, title, content, created_at, updated_at, tags
FROM posts
WHERE id = $1
`
	var post PostModel
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserId,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err

		}

	}
	return &post, nil
}
func (s *PostsStore) DeleteByID(ctx context.Context, id int64) (bool, error) {
	query := `
DELETE 
FROM posts
WHERE id = $1
RETURNING id
`
	var result int64
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&result,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return false, ErrNotFound
		default:
			return false, err

		}

	}
	return result == id, nil

}
