CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE teams (
    id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(255) NOT NULL,
    updated_at timestamp NOT NULL,
    created_at timestamp NOT NULL,
    deleted_at timestamp
);
