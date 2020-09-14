CREATE TABLE lists (
    id uuid PRIMARY KEY,
    date_created bigint,
    date_modified bigint,
    title text,
    description text,
    user_id uuid
);

CREATE TABLE items (
    id uuid PRIMARY KEY,
    date_created bigint,
    date_modified bigint,
    title text,
    description text,
    list_id uuid REFERENCES lists(id)
);

CREATE TABLE history (
    id uuid PRIMARY KEY,
    command text,
    state bytea,
    created bigint,
    user_id uuid
);
