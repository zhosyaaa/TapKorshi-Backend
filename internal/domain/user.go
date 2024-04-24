package domain

import "time"

//type User struct {
//	ID            uint
//	Email         string
//	Username      string
//	Phone         string
//	Password_hash string
//	Created_at    time.Time
//	LastVisitAt   time.Time
//	Verification  string
//	Posts         []Posts
//	Photo         string
//	Gender        string
//	Price         float64
//	Сommunal      bool
//	CommunalPrice float64
//	Contact       []string
//}

type User struct {
	ID            uint         `json:"ID,omitempty"`
	Email         string       `json:"email,omitempty"`
	Username      string       `json:"username,omitempty"`
	Phone         string       `json:"phone,omitempty"`
	Password_hash string       `json:"password_Hash,omitempty"`
	Created_at    time.Time    `json:"created_At"`
	LastVisitAt   time.Time    `json:"lastVisitAt"`
	Verification  Verification `json:"verification,omitempty"`
}
type Verification struct {
	Code     string `json:"code" bson:"code"`
	Verified bool   `json:"verified" bson:"verified"`
}
