package domain

import "time"

type Posts struct {
	ID      uint
	User    User
	Title   string
	Content string
	Type    string

	Created_at time.Time
	Location   Location
	Visibility string
	Tags       []string
	Images     []string
	Likes      int
	Comments   []Comment
}
