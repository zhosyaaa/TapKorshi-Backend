package repository

import (
	"database/sql"
	"github.com/zhosyaaa/RoommateTap/internal/domain"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user domain.User) error {
	query := `
        INSERT INTO users (
            email, username, phone, password_hash, created_at, last_visit_at, verification
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.db.Exec(
		query,
		user.Email, user.Username, user.Phone, user.Password_hash,
		user.Created_at, user.LastVisitAt, user.Verification,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Update(user domain.User) error {
	query := `
		UPDATE users 
		SET email = $1, username = $2, phone = $3, password_hash = $4, last_visit_at = $5, verification = $6 
		WHERE id = $7
	`
	_, err := r.db.Exec(
		query,
		user.Email, user.Username, user.Phone, user.Password_hash,
		user.LastVisitAt, user.Verification, user.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(userID uint) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByRefreshToken(refreshToken string) (domain.User, error) {
	var user domain.User
	query := `
		SELECT id, email, username, phone, password_hash, created_at, last_visit_at, verification 
		FROM users 
		WHERE refresh_token = $1 AND expires_at > $2
	`
	err := r.db.QueryRow(query, refreshToken, time.Now()).Scan(
		&user.ID, &user.Email, &user.Username, &user.Phone,
		&user.Password_hash, &user.Created_at, &user.LastVisitAt, &user.Verification,
	)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetByCredentials(email, password string) (domain.User, error) {
	var user domain.User
	query := `
		SELECT id, email, username, phone, password_hash, created_at, last_visit_at, verification 
		FROM users 
		WHERE email = $1 AND password_hash = $2
	`
	err := r.db.QueryRow(query, email, password).Scan(
		&user.ID, &user.Email, &user.Username, &user.Phone,
		&user.Password_hash, &user.Created_at, &user.LastVisitAt, &user.Verification,
	)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
