package services

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService interface {
	SendConfirmCode(toEmail string, toName string, confirmCode string) error
}

type emailService struct {
	apiKey     string
	templateID string
}

func NewEmailService() EmailService {
	return &emailService{
		apiKey:     os.Getenv("SENDGRID_API_KEY"),
		templateID: os.Getenv("CONFIRM_CODE_TEMPLATE_ID"),
	}
}

func (s *emailService) SendConfirmCode(toEmail string, toName string, confirmCode string) error {
	if s.apiKey == "" || s.templateID == "" {
		return fmt.Errorf("sendgrid not configured: API_KEY=%v, TEMPLATE_ID=%v", s.apiKey != "", s.templateID != "")
	}
	
	from := mail.NewEmail("Virtual Cuppa", "noreply@notacv.com")
	to := mail.NewEmail(toName, toEmail)
	
	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.SetTemplateID(s.templateID)
	
	personalization := mail.NewPersonalization()
	personalization.AddTos(to)
	personalization.SetDynamicTemplateData("Code", confirmCode)
	
	message.AddPersonalizations(personalization)
	
	client := sendgrid.NewSendClient(s.apiKey)
	response, err := client.Send(message)
	
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", toEmail, err)
	}
	
	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error for %s: status code %d, body: %s", toEmail, response.StatusCode, response.Body)
	}
	
	return nil
}
