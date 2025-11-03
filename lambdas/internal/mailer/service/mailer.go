package service

type Mailer interface {
	SendMail(to []string, subject string, body string) error
}
