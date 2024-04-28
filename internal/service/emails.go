package service

import (
	"fmt"
	"github.com/zhosyaaa/RoommateTap/internal/config"
	"github.com/zhosyaaa/RoommateTap/pkg/cache"
	"github.com/zhosyaaa/RoommateTap/pkg/email"
	"github.com/zhosyaaa/RoommateTap/pkg/email/sendpulse"
)

const (
	verificationLinkTmpl = "http://%s/api/v1/users/verify/%s" // https://<school host>/verification?code=<verification_code>
)

type EmailService struct {
	sender           email.Sender
	config           config.EmailConfig
	cache            cache.Cache
	sendpulseClients map[uint]*sendpulse.Client
}

func NewEmailsService(sender email.Sender, config config.EmailConfig, cache cache.Cache) *EmailService {
	return &EmailService{
		sender:           sender,
		config:           config,
		cache:            cache,
		sendpulseClients: make(map[uint]*sendpulse.Client),
	}
}

// Structures used for templates.
type verificationEmailInput struct {
	VerificationLink string
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
