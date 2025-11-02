package service

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type SesMailService struct {
	client    *sesv2.Client
	fromEmail string
}

func NewSesMailService(client *sesv2.Client, fromEmail string) *SesMailService {
	return &SesMailService{
		client:    client,
		fromEmail: strings.TrimSpace(fromEmail),
	}
}

func (sm *SesMailService) SendMail(to []string, subject string, body string) error {
	from := strings.TrimSpace(sm.fromEmail)
	if from == "" {
		return fmt.Errorf("from email is empty")
	}
	if _, err := mail.ParseAddress(from); err != nil {
		return fmt.Errorf("invalid from email: %w", err)
	}

	recips := make([]string, 0, len(to))
	for i := range to {
		addr := strings.TrimSpace(to[i])
		a, err := mail.ParseAddress(addr)
		if err != nil {
			return fmt.Errorf("invalid recipient: %w", err)
		}
		recips = append(recips, a.Address)
	}

	resp, err := sm.client.SendEmail(context.TODO(), &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(from),
		Destination: &types.Destination{
			ToAddresses: recips,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(subject)},
				Body: &types.Body{
					Text: &types.Content{Data: aws.String(body)},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error sending email: %v\n", err)
		return err
	}

	log.Printf("Email sent successfully, message ID: %s\n", aws.ToString(resp.MessageId))
	return nil
}
