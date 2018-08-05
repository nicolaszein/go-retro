CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE retrospectives (
    id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(255) NOT NULL,
    team_id uuid NOT NULL REFERENCES teams (id),
    updated_at timestamp NOT NULL,
    created_at timestamp NOT NULL,
    deleted_at timestamp
);
