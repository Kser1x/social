package store

import (
	"context"
	"database/sql"
	"errors"
)

type UserModel struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *UserModel) error {

	query := `INSERT INTO users(username,password,email) VALUES($1,$2,$3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
func (s *UsersStore) GetByID(ctx context.Context, userID int64) (*UserModel, error) {
	query := `
SELECT id, username, email, password, created_at
FROM users
WHERE id = $1
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &UserModel{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
		default:
			return nil, err

		}
	}
	return user, nil
}
func (s *UsersStore) DeleteByID(ctx context.Context, userID int64) (bool, error) {
	query := `
DELETE 
FROM users
WHERE id = $1
RETURNING id
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var result int64
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
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
	return result == userID, nil
}
func (s *UsersStore) Update(ctx context.Context, user *UserModel) error {
	query := `
UPDATE users
SET username = $1, email = $2
WHERE id = $3
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.ID)

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
