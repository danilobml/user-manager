package service

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type SesMailService struct {
	client  *sesv2.Client
	fromEmail string
}

func NewSesMailService(client  *sesv2.Client, fromEmail string) *SesMailService {
	return &SesMailService{
		client: client,
	}
}

func (sm *SesMailService) SendMail(to []string, subject string, body string) error {
	resp, err := sm.client.SendEmail(context.TODO(), &sesv2.SendEmailInput{
        FromEmailAddress: aws.String(sm.fromEmail),
        Destination: &types.Destination{
            ToAddresses: to,
        },
        Content: &types.EmailContent{
            Simple: &types.Message{
                Subject: &types.Content{
                    Data: aws.String(subject),
                },
                Body: &types.Body{
                    Text: &types.Content{
                        Data: aws.String(body),
                    },
                },
            },
        },
    })
    if err != nil {
        log.Printf("Error sending email: %v\n", err)
		return err
    }

    fmt.Printf("Email sent successfully, message ID: %s\n", *resp.MessageId)

	return nil
}

