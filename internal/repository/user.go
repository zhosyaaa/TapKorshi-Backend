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
func (r *UserRepository) SetSession(userID uint, session domain.Session) error {
	return nil
}

func (r *UserRepository) Create(user domain.User) (domain.User, error) {
	query := `
        INSERT INTO users (
            email, username, phone, password_hash, created_at, last_visit_at, verification_code, verification_verified
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id
    `
	var userID uint
	err := r.db.QueryRow(
		query,
		user.Email, user.Username, user.Phone, user.Password_hash,
		user.Created_at, user.LastVisitAt, user.VerificationCode, user.VerificationVerified,
	).Scan(&userID)
	if err != nil {
		return domain.User{}, err
	}

	user.ID = userID
	return user, nil
}

func (r *UserRepository) Update(user domain.User) error {
	query := `
		UPDATE users 
		SET email = $1, username = $2, phone = $3, password_hash = $4, last_visit_at = $5, verification_code = $6, verification_verified = $7
		WHERE id = $8
	`
	_, err := r.db.Exec(
		query,
		user.Email, user.Username, user.Phone, user.Password_hash,
		user.LastVisitAt, user.VerificationCode, user.VerificationVerified, user.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(userId uint) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByRefreshToken(refreshToken string) (domain.User, error) {
	var user domain.User
	query := `
		SELECT id, email, username, phone, password_hash, created_at, last_visit_at, verification_code, verification_verified
		FROM users 
		WHERE refresh_token = $1 AND expires_at > $2
	`
	err := r.db.QueryRow(query, refreshToken, time.Now()).Scan(
		&user.ID, &user.Email, &user.Username, &user.Phone,
		&user.Password_hash, &user.Created_at, &user.LastVisitAt, &user.VerificationCode, &user.VerificationVerified,
	)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetByCredentials(email, password string) (domain.User, error) {
	var user domain.User
	query := `
		SELECT id, email, username, phone, password_hash, created_at, last_visit_at, verification_code, verification_verified
		FROM users 
		WHERE email = $1 AND password_hash = $2
	`
	err := r.db.QueryRow(query, email, password).Scan(
		&user.ID, &user.Email, &user.Username, &user.Phone,
		&user.Password_hash, &user.Created_at, &user.LastVisitAt, &user.VerificationCode, &user.VerificationVerified,
	)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
func (r *UserRepository) Verify(userID uint, code string) error {
	query := `
		UPDATE users 
		SET verification_verified = true, verification_code = ''
		WHERE id = $1 AND verification_code = $2
	`
	res, err := r.db.Exec(query, userID, code)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrVerificationCodeInvalid
	}

	return nil
}
