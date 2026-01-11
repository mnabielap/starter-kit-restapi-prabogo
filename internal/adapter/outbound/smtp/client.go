package smtp_outbound_adapter

import (
	"fmt"
	"net/smtp"
	"os"

	outbound_port "prabogo/internal/port/outbound"
)

type smtpAdapter struct{}

func NewAdapter() outbound_port.EmailPort {
	return &smtpAdapter{}
}

func (s *smtpAdapter) SendEmail(to, subject, body string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("EMAIL_FROM")

	if os.Getenv("APP_MODE") == "test" {
		return nil // Skip sending in test mode
	}

	auth := smtp.PlainAuth("", username, password, host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=\"utf-8\"\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	addr := fmt.Sprintf("%s:%s", host, port)

	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}