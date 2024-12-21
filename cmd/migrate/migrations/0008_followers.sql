-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS followers(
user_id bigint NOT NULL,
followers_id bigint NOT NULL,
created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),


CONSTRAINT user_id_follower_id PRIMARY KEY (user_id, followers_id),
CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
CONSTRAINT fk_follower_id FOREIGN KEY (followers_id) REFERENCES users (id) ON DELETE CASCADE

    )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE followers;
-- +goose StatementEnd