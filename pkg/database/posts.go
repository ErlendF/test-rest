package database

import (
	"database/sql"
	"fmt"
	"test/pkg/models"
)

//GetPosts gets posts
func (db *Database) GetPosts() ([]models.Post, error) {
	var posts []models.Post
	stmt, err := db.Preparex(`SELECT id, content FROM posts;`)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("No rows")
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
		err = rows.Scan(&post)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	for _, p := range posts {
		p.Comments, err = db.getComments(p.ID)
	}

	return posts, nil
}

func (db *Database) getComments(ID int64) ([]models.Comment, error) {
	var comments []models.Comment
	stmt, err := db.Preparex(`SELECT id, content FROM comments WHERE post = $1;`)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("No rows")
		}

		return nil, err
	}

	rows, err := stmt.Queryx(ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comment models.Comment
	for rows.Next() {
		err = rows.Scan(&comment)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

//AddPost adds a new post
func (db *Database) AddPost(content string) error {
	result, err := db.Exec(`INSERT INTO posts (content) VALUES ($1);`, content)
	if err != nil {
		return err
	}
	num, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return fmt.Errorf("Could not insert post")
	}

	return nil
}

//AddComment adds a new comment to a post
func (db *Database) AddComment(comment *models.Comment) error {
	result, err := db.NamedExec(`INSERT INTO comment (post, content) VALUES (:post, :content);`, comment)
	if err != nil {
		return err
	}
	num, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return fmt.Errorf("Could not insert post")
	}

	return nil
}