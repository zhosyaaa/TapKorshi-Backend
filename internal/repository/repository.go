package repository

import (
	"github.com/zhosyaaa/RoommateTap/internal/domain"
)

type User interface {
	Create(user domain.User) error
	Update(user domain.User) error
	Delete(userid uint) error
	GetByRefreshToken(refreshToken string) (domain.User, error)
	GetByCredentials(email, password string) (user domain.User, err error)
	SetSession(userID uint, session domain.Session) error
}
