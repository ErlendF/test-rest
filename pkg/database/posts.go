package database

import (
	"database/sql"
	"errors"
	"test/pkg/models"
)

// GetPosts gets posts
func (db *Database) GetPosts() ([]models.Post, error) {
	var posts []models.Post
	stmt, err := db.Preparex(`SELECT id, content FROM posts;`)
	if err != nil {
		if err == sql.ErrNoRows {
			err = models.ErrNotFound
		}

		return nil, err
	}

	rows, err := stmt.Queryx()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var post models.Post
	for rows.Next() {
		err = rows.StructScan(&post)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	// getting all comments for each of the posts
	for i := range posts {
		posts[i].Comments, err = db.getComments(posts[i].ID)
		if err != nil && errors.Is(err, models.ErrNotFound) {
			return nil, err
		}
	}

	return posts, nil
}

func (db *Database) getComments(id int64) ([]models.Comment, error) {
	var comments []models.Comment
	stmt, err := db.Preparex(db.Rebind(`SELECT id, content FROM comments WHERE post = ?;`))
	if err != nil {
		if err == sql.ErrNoRows {
			err = models.ErrNotFound
		}

		return nil, err
	}

	rows, err := stmt.Queryx(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comment models.Comment
	for rows.Next() {
		err = rows.StructScan(&comment)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// AddPost adds a new post
func (db *Database) AddPost(content string) error {
	result, err := db.Exec(db.Rebind(`INSERT INTO posts (content) VALUES (?);`), content)
	if err != nil {
		return err
	}
	num, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return models.ErrDBInsert
	}

	return nil
}

// AddComment adds a new comment to a post
func (db *Database) AddComment(comment *models.Comment) error {
	result, err := db.NamedExec(`INSERT INTO comments (post, content) VALUES (:post, :content);`, comment)
	if err != nil {
		return err
	}
	num, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return models.ErrDBInsert
	}

	return nil
}
