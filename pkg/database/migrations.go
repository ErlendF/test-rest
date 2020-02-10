package database

import migrate "github.com/rubenv/sql-migrate"

// see gobuffalo/packr or markbates/pkger for alternatives to memorymigrations
func getMigrations() *migrate.MemoryMigrationSource {
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
