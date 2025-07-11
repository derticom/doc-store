-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE documents (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    mime TEXT NOT NULL,
    file BOOLEAN NOT NULL,
    public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
    grant_ids UUID[] DEFAULT '{}',
    json_data JSONB
);

-- +goose Down
DROP TABLE documents;
DROP TABLE users;
