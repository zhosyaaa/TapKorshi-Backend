package service

import (
	"fmt"
	"github.com/zhosyaaa/RoommateTap/internal/config"
	"github.com/zhosyaaa/RoommateTap/pkg/cache"
	"github.com/zhosyaaa/RoommateTap/pkg/email"
)

const (
	verificationLinkTmpl = "https://%s/verification?code=%s" // https://<school host>/verification?code=<verification_code>
)

type EmailService struct {
	sender email.Sender
	config config.EmailConfig
	cache  cache.Cache
}

// Structures used for templates.
type verificationEmailInput struct {
	VerificationLink string
}

type purchaseSuccessfulEmailInput struct {
	Name       string
	CourseName string
}

func (s *EmailService) SendUserVerificationEmail(input VerificationEmailInput) error {
	subject := fmt.Sprintf(s.config.Subjects.Verification, input.Name)

	templateInput := verificationEmailInput{s.createVerificationLink(input.Domain, input.VerificationCode)}
	sendInput := email.SendEmailInput{Subject: subject, To: input.Email}

	if err := sendInput.GenerateBodyFromHTML(s.config.Templates.Verification, templateInput); err != nil {
		return err
	}

	return s.sender.Send(sendInput)
}
func (s *EmailService) createVerificationLink(domain, code string) string {
	return fmt.Sprintf(verificationLinkTmpl, domain, code)
}
