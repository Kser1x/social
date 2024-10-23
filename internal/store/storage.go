package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *PostModel) error
	}
	Users interface {
		Create(ctx context.Context, user *UserModel) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostsStore{db},
		Users: &UsersStore{db},
	}
}
