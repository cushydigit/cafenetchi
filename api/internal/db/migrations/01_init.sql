
-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    phone TEXT UNIQUE NOT NULL,
    full_name TEXT,
    avatar_url TEXT,
    is_verified BOOLEAN DEFAULT false,
    status TEXT DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE user_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(id),
    role TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    agent_level TEXT,
    hourly_rate NUMERIC,
    rating NUMERIC DEFAULT 0,
    total_tasks INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE user_profiles;
DROP TABLE users;