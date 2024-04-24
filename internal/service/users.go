package service

import (
	"github.com/zhosyaaa/RoommateTap/internal/repository"
	"github.com/zhosyaaa/RoommateTap/pkg/auth"
	"github.com/zhosyaaa/RoommateTap/pkg/hash"
	"github.com/zhosyaaa/RoommateTap/pkg/otp"
)

type UserService struct {
	repo         repository.UserRepository
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	otpGenerator otp.Generator
	emailService Emails
}
