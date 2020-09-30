CREATE TABLE users (
    id uuid PRIMARY KEY,
    email text,
    passhash bytea,
    created bigint,
    modified bigint
);

CREATE TABLE version (
    version int PRIMARY KEY,
    created bigint
);

INSERT INTO version (version, created)
    SELECT 1, extract(epoch from now());
