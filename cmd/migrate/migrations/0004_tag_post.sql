-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    posts
    ADD
        COLUMN  tags VARCHAR(100) [];
ALTER TABLE posts ADD COLUMN updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE
post DROP COLUMN tags;
ALTER TABLE
post DROP COLUMN updated_at;
-- +goose StatementEnd
