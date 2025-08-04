package mailtrap

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Robotech-Org/gordian" // Your root package
	mail "github.com/xhit/go-simple-mail/v2"
)

// Config holds the necessary credentials for the Mailtrap service.
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	FromAddr string // e.g., "no-reply@diagramly.com"
}

// Emailer is a Mailtrap implementation of the gordian.Emailer interface.
type Emailer struct {
	config Config
}

// NewEmailer creates a new Mailtrap emailer.
func NewEmailer(cfg Config) *Emailer {
	return &Emailer{config: cfg}
}

// SendInvitation satisfies the interface.
func (e *Emailer) SendInvitation(ctx context.Context, invite *gordian.Invite) error {
	// --- Connect to the SMTP Server ---
	server := mail.NewSMTPClient()
	server.Host = e.config.Host
	port, _ := strconv.Atoi(e.config.Port)
	server.Port = port
	server.Username = e.config.Username
	server.Password = e.config.Password
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to mailtrap: %w", err)
	}

	// --- Create the Email ---
	email := mail.NewMSG()
	email.SetFrom(e.config.FromAddr).
		AddTo(invite.InviteeEmail).
		SetSubject("You're invited to join an organization!")

	// Create a simple invitation link for the email body
	invitationLink := fmt.Sprintf("https://app.diagramly.com/accept-invite?token=%s", invite.Token)

	// Set the email body (both HTML and plain text for compatibility)
	body := fmt.Sprintf("Hello! You have been invited to join an organization. Please click the link to accept: %s", invitationLink)
	htmlBody := fmt.Sprintf("<h1>Hello!</h1><p>You have been invited to join an organization. Please click the link below to accept:</p><a href=\"%s\">Accept Invitation</a>", invitationLink)
	email.SetBody(mail.TextPlain, body)
	email.AddAlternative(mail.TextHTML, htmlBody)

	// --- Send the Email ---
	err = email.Send(smtpClient)
	if err != nil {
		return fmt.Errorf("failed to send email via mailtrap: %w", err)
	}

	return nil
}
