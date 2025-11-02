package mocks

type MockMailer struct {
	To      []string
	Subject string
	Message string
}

func (m *MockMailer) SendMail(to []string, subject, body string) error {
	m.To = append([]string(nil), to...)
	m.Subject = subject
	m.Message = "Subject: " + subject + "\n" + body
	return nil
}