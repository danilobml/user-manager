package service

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/danilobml/user-manager/internal/errs"
)

type LocalMailConfig struct {
	FromEmail     string
	FromEmailPass string
	FromEmailSMTP string
	SMTPAddr      string
}

type LocalMailService struct {
	fromEmail     string
	fromEmailPass string
	fromEmailSMTP string
	smtpAddr      string
}

func NewLocalMailService(cfg LocalMailConfig) *LocalMailService {
	return &LocalMailService{
		fromEmail:     cfg.FromEmail,
		fromEmailPass: cfg.FromEmailPass,
		fromEmailSMTP: cfg.FromEmailSMTP,
		smtpAddr:      cfg.SMTPAddr,
	}
}

func (ms *LocalMailService) SendMail(to []string, subject, body string) error {
	if ms.fromEmail == "" || ms.fromEmailPass == "" || ms.fromEmailSMTP == "" || ms.smtpAddr == "" {
		return errs.ErrMailServiceDisabled
	}

	var msg bytes.Buffer
	msg.WriteString("From: " + ms.fromEmail + "\r\n")
	msg.WriteString("To: " + strings.Join(to, ", ") + "\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	msg.WriteString("Content-Transfer-Encoding: 7bit\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body + "\r\n")

	auth := smtp.PlainAuth("", ms.fromEmail, ms.fromEmailPass, ms.fromEmailSMTP)

	fmt.Printf("Email sent successfully to: %s\n", to[0])
	return smtp.SendMail(ms.smtpAddr, auth, ms.fromEmail, to, msg.Bytes())
}
