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
    created bigint,
    user_id uuid
);
