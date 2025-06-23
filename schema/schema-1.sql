CREATE TABLE IF NOT EXISTS users(
    name text,
    email text primary key,
    password text
);

CREATE TABLE IF NOT EXISTS todo(
    id serial primary key,
    title text,
    description text,
    email text
);