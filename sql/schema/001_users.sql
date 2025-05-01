-- +goose Up 
CREATE TABLE users( 
	id int NOT NULL PRIMARY KEY,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	email TEXT
);


-- +goose Down
DROP TABLE users;
