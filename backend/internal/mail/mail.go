package mail

import (
	"context"
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Mailer interface {
	SendVerificationEmail(ctx context.Context, toEmail, otp, username string) error
}

type sendgridMailer struct {
	client    *sendgrid.Client
	fromEmail string
	fromName  string
}

func NewSendGridMailer() Mailer {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	fromEmail := os.Getenv("SENDGRID_FROM_EMAIL")
	fromName := os.Getenv("SENDGRID_FROM_NAME")

	return &sendgridMailer{
		client:    sendgrid.NewSendClient(apiKey),
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

func (m *sendgridMailer) SendVerificationEmail(ctx context.Context, toEmail, otp, username string) error {
	from := mail.NewEmail(m.fromName, m.fromEmail)
	subject := fmt.Sprintf("%s is your verification code", otp) // Putting OTP in the subject line increases open-rates
	to := mail.NewEmail(username, toEmail)

	// Clean, scannable transactional email layout focused entirely on the OTP code
	htmlContent := fmt.Sprintf(`
		<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e1e1e1; border-radius: 8px;">
			<h2 style="color: #333;">Welcome to Relationship Memory Platform, %s!</h2>
			<p style="color: #555; font-size: 16px;">Please use the following One-Time Password (OTP) to verify your email address:</p>
			
			<div style="margin: 30px 0; text-align: center;">
				<div style="display: inline-block; padding: 15px 30px; background-color: #F3F4F6; border: 2px dashed #4F46E5; font-size: 32px; font-weight: bold; letter-spacing: 5px; color: #4F46E5; border-radius: 6px;">
					%s
				</div>
			</div>
			
			<p style="color: #9CA3AF; font-size: 14px;">This code is private and will expire in 15 minutes. If you did not request this code, please ignore this email.</p>
		</div>
	`, username, otp)

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	response, err := m.client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid failed with status code: %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}
