package model

type Mailer interface {
	SendMail(to []string, subject string, body string) error
}
