CREATE TABLE users (
    id uuid PRIMARY KEY,
    email text,
    passhash bytea,
    created bigint,
    modified bigint
);
