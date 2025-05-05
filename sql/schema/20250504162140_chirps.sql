-- +goose Up
-- +goose StatementBegin
CREATE TABLE chirps( 
	id UUID NOT NULL PRIMARY KEY,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	body TEXT NOT NULL,
	user_id UUID NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chirps;
-- +goose StatementEnd
