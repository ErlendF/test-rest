package models

//Post is a post and related comments
type Post struct {
	ID       int64     `json:"id"`
	Content  string    `json:"content"`
	Comments []Comment `json:"comments"`
}
