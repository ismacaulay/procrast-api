CREATE TABLE lists (
    id uuid PRIMARY KEY,
    title text,
    description text,
    created bigint,
    modified bigint,
    user_id uuid
);

CREATE TABLE items (
    id uuid PRIMARY KEY,
    title text,
    description text,
    state smallint,
    created bigint,
    modified bigint,
    list_id uuid REFERENCES lists(id)
);

CREATE TABLE history (
    id uuid PRIMARY KEY,
    command text,
    state bytea,
    ts bigint,
    created bigint,
    user_id uuid
);

CREATE TABLE version (
    version int PRIMARY KEY,
    created bigint
);

INSERT INTO version (version, created)
    SELECT 1, extract(epoch from now());
