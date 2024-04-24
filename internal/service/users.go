package service

import (
	"context"
	"errors"
	"github.com/zhosyaaa/RoommateTap/internal/domain"
	"github.com/zhosyaaa/RoommateTap/internal/repository"
	"github.com/zhosyaaa/RoommateTap/pkg/auth"
	"github.com/zhosyaaa/RoommateTap/pkg/hash"
	"github.com/zhosyaaa/RoommateTap/pkg/otp"
	"time"
)

type UsersService struct {
	repo         repository.UserRepository
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	otpGenerator otp.Generator
	emailService Emails

	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
	verificationCodeLength int

	domain string
}

func NewUsersService(repo repository.UserRepository, hasher hash.PasswordHasher, tokenManager auth.TokenManager, otpGenerator otp.Generator, emailService Emails, accessTokenTTL time.Duration, refreshTokenTTL time.Duration, verificationCodeLength int, domain string) *UsersService {
	return &UsersService{repo: repo, hasher: hasher, tokenManager: tokenManager, otpGenerator: otpGenerator, emailService: emailService, accessTokenTTL: accessTokenTTL, refreshTokenTTL: refreshTokenTTL, verificationCodeLength: verificationCodeLength, domain: domain}
}

func (s *UsersService) SignUp(ctx context.Context, input UserSignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}
	verificationCode := s.otpGenerator.RandomSecret(s.verificationCodeLength)
	user := domain.User{
		Username:      input.Username,
		Password_hash: passwordHash,
		Phone:         input.Phone,
		Email:         input.Email,
		Created_at:    time.Now(),
		LastVisitAt:   time.Now(),
		Verification: domain.Verification{
			Code: verificationCode,
		},
	}
	if err := s.repo.Create(user); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return err
		}

		return err
	}
	return s.emailService.SendUserVerificationEmail(VerificationEmailInput{
		Email:            user.Email,
		Name:             user.Username,
		VerificationCode: verificationCode,
	})
}
func (s *UsersService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}
	user, err := s.repo.GetByCredentials(input.Email, passwordHash)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return Tokens{}, err
		}

		return Tokens{}, err
	}
	return s.createSession(ctx, user.ID)
}
func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	return Tokens{}, nil
}

func (s *UsersService) createSession(ctx context.Context, userId primitive.ObjectID) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId.Hex(), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, userId, session)

	return res, err
}
