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
	currentKey := os.Getenv("SENDGRID_API_KEY")

	fmt.Println("\n--- [MAILER DEBUG TRIGGERED] ---")
	if currentKey == "" {
		fmt.Println("[ERROR] SENDGRID_API_KEY is completely EMPTY at runtime!")
	} else if len(currentKey) < 8 {
		fmt.Printf("[WARNING] SendGrid key is active but suspiciously short: '%s'\n", currentKey)
	} else {
		fmt.Printf("[OK] Runtime SendGrid API Key starts with: %s...\n", currentKey[:8])
	}
	fmt.Printf("[OK] From Email Configured As: %s\n", m.fromEmail)
	fmt.Println("--------------------------------")

	from := mail.NewEmail(m.fromName, m.fromEmail)
	subject := "Verify your RealtimeChatty Account"
	to := mail.NewEmail(username, toEmail)

	// Clean plain text fallback for watch screens / pure-text mail readers
	plainTextContent := fmt.Sprintf("Hello %s, welcome to RealtimeChatty! Your 6-digit verification code is: %s. This code expires in 15 minutes.", username, otp)

	// Get our styled HTML payload
	htmlContent := getHTMLTemplate(username, otp)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	response, err := m.client.SendWithContext(ctx, message)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid failed with status code: %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}

// Helper to bundle our responsive custom styled CSS layout cleanly
func getHTMLTemplate(username, otp string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify your Email</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
            background-color: #f4f6f8;
            margin: 0;
            padding: 0;
            -webkit-font-smoothing: antialiased;
        }
        .email-container {
            max-width: 550px;
            margin: 40px auto;
            background-color: #ffffff;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
            border: 1px solid #eef2f5;
        }
        .header {
            background-color: #4f46e5; /* Modern Indigo */
            padding: 32px;
            text-align: center;
        }
        .header h1 {
            color: #ffffff;
            margin: 0;
            font-size: 24px;
            font-weight: 700;
            letter-spacing: -0.5px;
        }
        .body-content {
            padding: 40px 32px;
            color: #334155;
            line-height: 1.6;
        }
        .body-content h2 {
            font-size: 20px;
            color: #1e293b;
            margin-top: 0;
            margin-bottom: 16px;
        }
        .body-content p {
            margin-top: 0;
            margin-bottom: 24px;
            font-size: 15px;
            color: #64748b;
        }
        .otp-container {
            background-color: #f8fafc;
            border: 2px dashed #cbd5e1;
            border-radius: 8px;
            padding: 20px;
            text-align: center;
            margin: 28px 0;
        }
        .otp-code {
            font-size: 36px;
            font-weight: 800;
            color: #4f46e5;
            letter-spacing: 6px;
            margin: 0;
            font-family: 'Courier New', Courier, monospace;
        }
        .footer {
            background-color: #f8fafc;
            padding: 24px 32px;
            text-align: center;
            border-top: 1px solid #eef2f5;
        }
        .footer p {
            margin: 0;
            font-size: 12px;
            color: #94a3b8;
        }
        .warning {
            font-size: 13px !important;
            color: #94a3b8 !important;
            border-top: 1px solid #f1f5f9;
            padding-top: 16px;
            margin-bottom: 0 !important;
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <h1>RealtimeChatty</h1>
        </div>
        <div class="body-content">
            <h2>Hello, %s!</h2>
            <p>Thank you for signing up. To complete your account registration and access your dashboard, please confirm your email address by entering the verification security code below:</p>
            
            <div class="otp-container">
                <h3 class="otp-code">%s</h3>
            </div>
            
            <p>This verification code is sensitive and is valid for the next <strong>15 minutes</strong>. For your security, do not share this token with anyone.</p>
            <p class="warning">If you did not request this registration sequence, you can safely ignore this automated message—no further configuration is required.</p>
        </div>
        <div class="footer">
            <p>&copy; 2026 RealtimeChatty Platform. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, username, otp)
}
