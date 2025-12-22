package services

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService interface {
	SendConfirmCode(toEmail string, toName string, confirmCode string) error
	SendInvitation(toEmail string, toName string, organisationName string) error
	SendMatchAccepted(toEmail string, toName string, matchName string, availability string) error
}

type emailService struct {
	apiKey                  string
	confirmCodeTemplateID   string
	invitationTemplateID    string
	matchAcceptedTemplateID string
}

func NewEmailService() EmailService {
	return &emailService{
		apiKey:                  os.Getenv("SENDGRID_API_KEY"),
		confirmCodeTemplateID:   os.Getenv("CONFIRM_CODE_TEMPLATE_ID"),
		invitationTemplateID:    os.Getenv("INVITATION_TEMPLATE_ID"),
		matchAcceptedTemplateID: os.Getenv("MATCH_ACCEPTED_TEMPLATE_ID"),
	}
}

func (s *emailService) SendConfirmCode(toEmail string, toName string, confirmCode string) error {
	if s.apiKey == "" || s.confirmCodeTemplateID == "" {
		return fmt.Errorf("sendgrid not configured: API_KEY=%v, TEMPLATE_ID=%v", s.apiKey != "", s.confirmCodeTemplateID != "")
	}
	
	from := mail.NewEmail("Virtual Cuppa", "noreply@notacv.com")
	to := mail.NewEmail(toName, toEmail)
	
	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.SetTemplateID(s.confirmCodeTemplateID)
	
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

func (s *emailService) SendInvitation(toEmail string, toName string, organisationName string) error {
	if s.apiKey == "" || s.invitationTemplateID == "" {
		return fmt.Errorf("sendgrid not configured for invitations: API_KEY=%v, TEMPLATE_ID=%v", s.apiKey != "", s.invitationTemplateID != "")
	}
	
	from := mail.NewEmail("Virtual Cuppa", "noreply@notacv.com")
	to := mail.NewEmail(toName, toEmail)
	
	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.SetTemplateID(s.invitationTemplateID)
	
	personalization := mail.NewPersonalization()
	personalization.AddTos(to)
	personalization.SetDynamicTemplateData("OrganisationName", organisationName)
	
	message.AddPersonalizations(personalization)
	
	client := sendgrid.NewSendClient(s.apiKey)
	response, err := client.Send(message)
	
	if err != nil {
		return fmt.Errorf("failed to send invitation email to %s: %w", toEmail, err)
	}
	
	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error for invitation to %s: status code %d, body: %s", toEmail, response.StatusCode, response.Body)
	}
	
	return nil
}

func (s *emailService) SendMatchAccepted(toEmail string, toName string, matchName string, availability string) error {
	if s.apiKey == "" || s.matchAcceptedTemplateID == "" {
		return fmt.Errorf("sendgrid not configured for match accepted notifications: API_KEY=%v, TEMPLATE_ID=%v", s.apiKey != "", s.matchAcceptedTemplateID != "")
	}
	
	from := mail.NewEmail("Virtual Cuppa", "noreply@notacv.com")
	to := mail.NewEmail(toName, toEmail)
	
	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.SetTemplateID(s.matchAcceptedTemplateID)
	
	personalization := mail.NewPersonalization()
	personalization.AddTos(to)
	personalization.SetDynamicTemplateData("MatchName", matchName)
	personalization.SetDynamicTemplateData("Availability", availability)
	
	message.AddPersonalizations(personalization)
	
	client := sendgrid.NewSendClient(s.apiKey)
	response, err := client.Send(message)
	
	if err != nil {
		return fmt.Errorf("failed to send match accepted notification to %s: %w", toEmail, err)
	}
	
	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error for match accepted notification to %s: status code %d, body: %s", toEmail, response.StatusCode, response.Body)
	}
	
	return nil
}
