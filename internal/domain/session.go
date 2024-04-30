package domain

import "time"

type Session struct {
	Userid       uint      `json:"userid"`
	RefreshToken string    `json:"refreshToken"`
	Fingerprint  string    `json:"fingerprint"`
	Ip           string    `json:"ip"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
