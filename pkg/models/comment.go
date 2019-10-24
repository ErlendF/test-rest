package models

//Comment is a simple comment on a post
type Comment struct {
	ID      int64  `json:"id" db:"id"`
	Post    int64  `json:"-" db:"post"`
	Content string `json:"content" db:"content"`
}
