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
		Create(context.Context, *PostModel) error
		DeleteByID(context.Context, int64) (bool, error)
		Update(context.Context, *PostModel) error
		GetUserFeed(context.Context, int64) ([]PostWithMetadata, error)
	}
	Users interface {
		GetByID(context.Context, int64) (*UserModel, error)
		DeleteByID(context.Context, int64) (bool, error)
		Create(context.Context, *UserModel) error
		Update(context.Context, *UserModel) error
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}
	Followers interface {
		Follow(ctx context.Context, followerID, userID int64) error
		Unfollow(ctx context.Context, followerID, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostsStore{db},
		Users:     &UsersStore{db},
		Comments:  &CommentsStore{db},
		Followers: &FollowerStore{db},
	}
}
