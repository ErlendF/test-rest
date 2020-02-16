package database

import migrate "github.com/rubenv/sql-migrate"

// see gobuffalo/packr or markbates/pkger for alternatives to memorymigrations
func (db *Database) getMigrations() *migrate.MemoryMigrationSource {
	if db.dbType == "postgres" {
		return &migrate.MemoryMigrationSource{
			Migrations: []*migrate.Migration{
				{
					Id: "001",
					Up: []string{
						`CREATE TABLE "posts" (
							"id" bigserial PRIMARY KEY,
							"content" text NOT NULL
						);`,

						`CREATE TABLE "comments" (
							"id" bigserial PRIMARY KEY,
							"post" bigint REFERENCES "posts"(id) ON DELETE CASCADE NOT NULL,
							"content" text NOT NULL
						);`,
					},
					Down: []string{"DROP TABLE posts", "DROP TABLE comments"},
				},
			},
		}
	}

	return &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "001",
				Up: []string{
					`CREATE TABLE posts (
						id bigint NOT NULL AUTO_INCREMENT,
						content text NOT NULL,
						CONSTRAINT posts_pk PRIMARY KEY (id)
					);`,

					`CREATE TABLE comments (
						id bigint NOT NULL AUTO_INCREMENT,
						post bigint NOT NULL REFERENCES posts(id),
						content text NOT NULL,
						CONSTRAINT comments_pk PRIMARY KEY (id)
					);`,
				},
				Down: []string{"DROP TABLE posts", "DROP TABLE comments"},
			},
		},
	}
}
