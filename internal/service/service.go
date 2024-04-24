package service

import "context"

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
	SendStudentVerificationEmail(VerificationEmailInput) error
	SendUserVerificationEmail(VerificationEmailInput) error
}
