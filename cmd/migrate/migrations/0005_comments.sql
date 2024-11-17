-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS comments(
                                    id bigserial PRIMARY KEY,
                                    post_id bigserial NOT NULL,
                                    user_id bigserial NOT NULL,
                                    content text NOT NULL,
                                    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()

    )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE comments;
-- +goose StatementEnd