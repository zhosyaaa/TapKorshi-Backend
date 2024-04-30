package domain

import "time"

type Session struct {
	Userid    uint      `json:"userid"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expiresAt"`
}
