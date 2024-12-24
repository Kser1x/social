package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type PostModel struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserId    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      UserModel `json:"user"`
}

type PostWithMetadata struct {
	PostModel
	CommentCount int `json:"comment_count"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) GetUserFeed(ctx context.Context, userID int64) ([]PostWithMetadata, error) {
	query := `
SELECT
p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username,
COUNT(c.id) AS comments_count
FROM posts p 
LEFT JOIN comments c ON c.post_id = p.id
LEFT JOIN  users u ON p.user_id = u.id
JOIN  followers f ON  f.followers_id = p.user_id OR p.user_id = $1
WHERE  f.user_id = $1 OR  p.user_id = $1
GROUP BY p.id, u.username
ORDER BY p.created_at DESC 
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata
		err := rows.Scan(
			&p.ID,
			&p.UserId,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, p)

	}
	return feed, nil
}

func (s *PostsStore) Create(ctx context.Context, post *PostModel) error {
	query := `INSERT INTO posts (content, title, user_id, tags)
				VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
SELECT id, user_id, title, content, created_at, updated_at, tags, version
FROM posts
WHERE id = $1
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post PostModel
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserId,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

func (s *PostsStore) Update(ctx context.Context, post *PostModel) error {
	version := post.Version + 1
	query := `
UPDATE posts
SET title = $1, content = $2, version = $5
WHERE id = $3 AND version = $4
RETURNING version
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
		version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err

		}
	}
	return err
}
