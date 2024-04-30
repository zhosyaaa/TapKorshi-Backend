package repository

import (
	"database/sql"
	"github.com/zhosyaaa/RoommateTap/internal/domain"
)

type User interface {
	Create(user domain.User) (domain.User, error)
	Update(user domain.User) error
	Delete(userid uint) error
	GetByCredentials(email, password string) (user domain.User, err error)
	Verify(userID uint, code string) error
}

type Repositories struct {
	Users User
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users: NewUserRepository(db),
	}
}
