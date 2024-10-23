-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts(
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    user_id bigint NOT NULL,
    content text NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
    )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd