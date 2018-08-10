CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE cards (
    id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    content text NOT NULL,
    retrospective_id uuid NOT NULL REFERENCES retrospectives (id),
    votes int,
    type varchar(100) NOT NULL,
    updated_at timestamp NOT NULL,
    created_at timestamp NOT NULL,
    deleted_at timestamp
);
