-- the urls table
CREATE TABLE IF NOT EXISTS urls (
	id VARCHAR(12) UNIQUE,
	url TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- hash index on id
CREATE INDEX IF NOT EXISTS urls_id_hash_idx ON urls USING HASH (id);
