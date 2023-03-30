package services

import (
	"fmt"

	"github.com/avalonprod/gasstrem/src/internal/config"
	"github.com/avalonprod/gasstrem/src/packages/email"
)

type VerificationEmailInput struct {
	Email            string
	Name             string
	VerificationCode string
}

type EmailsService struct {
	sender email.Sender
	config config.EmailConfig
}

func NewEmailService(sender email.Sender, config config.EmailConfig) *EmailsService {
	return &EmailsService{
		sender: sender,
		config: config,
	}
}

type verificationEmailInput struct {
	Name             string
	VerificationCode string
}

func (e *EmailsService) SendUserVerificationEmail(input VerificationEmailInput) error {
	subject := fmt.Sprintf(e.config.Subjects.Verification, input.Name)
	sendInput := email.SendEmailInput{To: input.Email, Subject: subject}

	templateInput := verificationEmailInput{Name: input.Name, VerificationCode: input.VerificationCode}

	err := sendInput.GenerateBodyFromHTML(e.config.Templates.Verification, templateInput)
	if err != nil {
		return err
	}
	err = e.sender.Send(sendInput)
	return err
}
