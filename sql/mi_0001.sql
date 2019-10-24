-- +migrate Up

CREATE TABLE "posts" (
    "id" bigserial PRIMARY KEY,
    "content" text NOT NULL,
);

CREATE TABLE "comments" (
    "id" bigserial PRIMARY KEY,
    "post" bigint REFERENCES "posts"(id) ON DELETE CASCADE NOT NULL,
    "content" text NOT NULL,
);

-- +migrate