package service

import (
	"context"
	"github.com/zhosyaaa/RoommateTap/internal/config"
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

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	//Verify(ctx context.Context, userID primitive.ObjectID, hash string) error
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
	Users  Users
	Emails Emails
}

func NewServices(deps Deps) *Services {
	emailsService := NewEmailsService(deps.EmailSender, deps.EmailConfig, deps.Cache)
	usersService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.OtpGenerator, emailsService,
		deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.VerificationCodeLength, deps.Domain)
	return &Services{
		Users: usersService,
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
	//DNS                    dns.DomainManager
}
