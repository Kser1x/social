package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("record not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		GetByID(context.Context, int64) (*PostModel, error)
		Create(ctx context.Context, post *PostModel) error
		DeleteByID(context.Context, int64) (bool, error)
		Update(ctx context.Context, model *PostModel) error
	}
	Users interface {
		GetByID(context.Context, int64) (*UserModel, error)
		DeleteByID(context.Context, int64) (bool, error)
		Create(ctx context.Context, user *UserModel) error
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentsStore{db},
	}
}
