-- +goose Up 
CREATE TABLE users( 
	id UUID NOT NULL PRIMARY KEY,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	email TEXT,
	hashed_password TEXT NOT NULL	
);


-- +goose Down
DROP TABLE users;
