package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/zhosyaaa/RoommateTap/internal/config"
	"github.com/zhosyaaa/RoommateTap/internal/domain"
	"github.com/zhosyaaa/RoommateTap/internal/repository"
	"github.com/zhosyaaa/RoommateTap/pkg/auth"
	"github.com/zhosyaaa/RoommateTap/pkg/cache"
	"github.com/zhosyaaa/RoommateTap/pkg/email"
	"github.com/zhosyaaa/RoommateTap/pkg/hash"
	"github.com/zhosyaaa/RoommateTap/pkg/otp"
	"time"
)

type UserSignUpInput struct {
	Username string
	Email    string
	Phone    string
	Password string
}
type UserSignInInput struct {
	Email    string
	Password string
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}
type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	SignUp(ctx context.Context, input UserSignUpInput, Fingerprint, IP string) (Tokens, string, error)
	SignIn(ctx context.Context, input UserSignInInput, Fingerprint, IP string) (Tokens, string, error)
	Verify(ctx context.Context, userID uint, hash string) error
	RefreshTokens(sessionId, token, fingerprint string) (Tokens, string, error)
	OAuthSignIn(ctx context.Context, googleUser GoogleUser, Fingerprint, IP string) (Tokens, string, error)
}

type VerificationEmailInput struct {
	Email            string
	Name             string
	VerificationCode string
	Domain           string
}

type Emails interface {
	SendUserVerificationEmail(VerificationEmailInput) error
}

type Services struct {
	Users    Users
	Emails   Emails
	Sessions Sessions
}

type Sessions interface {
	CreateSession(session *domain.Session) (string, error)
	GetSession(sessionID string) (*domain.Session, error)
	DeleteSession(sessionID string) error
}

func NewServices(deps Deps) *Services {
	emailsService := NewEmailsService(deps.EmailSender, deps.EmailConfig, deps.Cache)
	sessionService := NewSessionService(deps.RedisClient)
	usersService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.OtpGenerator, emailsService, sessionService,
		deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.VerificationCodeLength, deps.Domain)
	return &Services{
		Users:    usersService,
		Emails:   emailsService,
		Sessions: sessionService,
	}
}

type Deps struct {
	Repos        *repository.Repositories
	Cache        cache.Cache
	Hasher       hash.PasswordHasher
	TokenManager auth.TokenManager
	EmailSender  email.Sender
	EmailConfig  config.EmailConfig
	//StorageProvider        storage.Provider
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	//FondyCallbackURL       string
	CacheTTL               int64
	OtpGenerator           otp.Generator
	VerificationCodeLength int
	Environment            string
	Domain                 string
	RedisClient            *redis.Client
	//DNS                    dns.DomainManager
}
