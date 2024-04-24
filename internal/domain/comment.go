package domain

import "time"

type Comment struct {
	User       User
	Content    string
	Created_at time.Time
}
