package models

import "errors"

// ErrNotFound indicates that a requested resource was not found
var ErrNotFound = errors.New("not found")

// ErrDBInsert indicates that there were some error inserting something in the database
var ErrDBInsert = errors.New("could not insert item into the database")
