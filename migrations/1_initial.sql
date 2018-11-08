-- +migrate Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS user_profile (
    user_id serial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    password varchar(64) NOT NULL,

    nickname citext UNIQUE NOT NULL,
    avatar text,

    record integer DEFAULT 0,
    win integer DEFAULT 0,
    draws integer DEFAULT 0,
    loss integer DEFAULT 0
);

-- +migrate Down
DROP TABLE IF EXISTS user_profile;

DROP EXTENSION IF EXISTS citext;
