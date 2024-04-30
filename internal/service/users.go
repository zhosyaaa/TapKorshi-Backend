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
	repo                   repository.User
	hasher                 hash.PasswordHasher
	tokenManager           auth.TokenManager
	otpGenerator           otp.Generator
	emailService           Emails
	sessionService         Sessions
	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
	verificationCodeLength int

	domain string
}

func (s *UsersService) Verify(ctx context.Context, userID uint, hash string) error {
	err := s.repo.Verify(userID, hash)
	if err != nil {
		if errors.Is(err, domain.ErrVerificationCodeInvalid) {
			return err
		}

		return err
	}

	return nil
}

func NewUsersService(repo repository.User, hasher hash.PasswordHasher, tokenManager auth.TokenManager, otpGenerator otp.Generator, emailService Emails, sessionService Sessions, accessTokenTTL time.Duration, refreshTokenTTL time.Duration, verificationCodeLength int, domain string) *UsersService {
	return &UsersService{repo: repo, hasher: hasher, tokenManager: tokenManager, otpGenerator: otpGenerator, emailService: emailService, sessionService: sessionService, accessTokenTTL: accessTokenTTL, refreshTokenTTL: refreshTokenTTL, verificationCodeLength: verificationCodeLength, domain: domain}
}

func (s *UsersService) SignUp(ctx context.Context, input UserSignUpInput, Fingerprint, IP string) (Tokens, string, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, "", err
	}
	verificationCode := s.otpGenerator.RandomSecret(s.verificationCodeLength)
	user := domain.User{
		Username:             input.Username,
		Password_hash:        passwordHash,
		Phone:                input.Phone,
		Email:                input.Email,
		Created_at:           time.Now(),
		LastVisitAt:          time.Now(),
		VerificationCode:     verificationCode,
		VerificationVerified: false,
	}
	if user, err = s.repo.Create(user); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return Tokens{}, "", err
		}
		return Tokens{}, "", err
	}
	s.emailService.SendUserVerificationEmail(VerificationEmailInput{
		Email:            user.Email,
		Name:             user.Username,
		VerificationCode: verificationCode,
		Domain:           "localhost:8000",
	})
	return s.createSession(user.ID, Fingerprint, IP)
}

func (s *UsersService) SignIn(ctx context.Context, input UserSignInInput, Fingerprint, IP string) (Tokens, string, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, "", err
	}
	user, err := s.repo.GetByCredentials(input.Email, passwordHash)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return Tokens{}, "", err
		}

		return Tokens{}, "", err
	}
	return s.createSession(user.ID, Fingerprint, IP)
}

func (s *UsersService) createSession(userId uint, fingerprint, ip string) (Tokens, string, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId, s.accessTokenTTL)
	if err != nil {
		return res, "", err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, "", err
	}
	session := &domain.Session{
		Userid:       userId,
		RefreshToken: res.RefreshToken,
		Fingerprint:  fingerprint,
		Ip:           ip,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}
	sessionId, err := s.sessionService.CreateSession(session)

	return res, sessionId, err
}
func (s *UsersService) RefreshTokens(sessionId, token, fingerprint string) (Tokens, string, error) {
	var (
		res       Tokens
		err       error
		sessionID string
	)
	session, err := s.sessionService.GetSession(sessionId)
	if err != nil {
		return res, sessionID, err
	}
	if time.Now().After(session.ExpiresAt) {
		return res, sessionID, errors.New("Refresh token expired")
	}
	if session.RefreshToken != token || session.Fingerprint != fingerprint {
		return res, sessionID, errors.New("invalid refresh token or fingerprint")
	}

	err = s.sessionService.DeleteSession(sessionId)
	if err != nil {
		return res, sessionID, err
	}
	return s.createSession(session.Userid, fingerprint, session.Ip)
}
