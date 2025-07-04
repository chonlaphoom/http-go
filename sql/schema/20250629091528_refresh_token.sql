-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens( 
	token TEXT NOT NULL PRIMARY KEY,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	expires_at TIMESTAMP, 
	revoked_at TIMESTAMP NULL,
	user_id UUID NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd
