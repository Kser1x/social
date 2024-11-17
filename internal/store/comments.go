package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"`
	User      UserModel `json:"user"`
}

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
  SELECT  c.content, c.created_at, u.id, u.username
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = $1
    ORDER BY c.created_at DESC;
`
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.Content, &c.CreatedAt, &c.User.ID, &c.User.Username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)

	}
	return comments, nil

}
