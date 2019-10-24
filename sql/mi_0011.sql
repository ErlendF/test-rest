-- +migrate Up

CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "username" text NOT NULL
);

CREATE TABLE "posts" (
    "id" bigserial PRIMARY KEY,
    "author" bigint REFERENCES "users"(id) ON DELETE CASCADE NOT NULL,
    "content" text NOT NULL,
);

CREATE TABLE "comments" (
    "id" bigserial PRIMARY KEY,
    "author" bigint REFERENCES "users"(id) ON DELETE CASCADE NOT NULL,
    "post"bigint REFERENCES "posts"(id) ON DELETE CASCADE NOT NULL,
    "content" text NOT NULL,
);

-- +migrate