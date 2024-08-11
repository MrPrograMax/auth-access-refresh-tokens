package email

import (
	"errors"
	"github.com/go-gomail/gomail"
)

type SMTPSender struct {
	from string
	pass string
	host string
	port int
}

func NewSMTPSender(from, pass, host string, port int) (*SMTPSender, error) {
	if !IsEmailValid(from) {
		return nil, errors.New("invalid \"from\" email")
	}

	return &SMTPSender{
		from: from, pass: pass, host: host, port: port}, nil
}

func (s *SMTPSender) Send(input SendEmailInput) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", input.To)
	m.SetHeader("Subject", input.Subject)
	m.SetBody("text/plain", input.Body)

	dialer := gomail.NewDialer(s.host, s.port, s.from, s.pass)
	if err := dialer.DialAndSend(m); err != nil {
		return errors.New("failed to send email via smtp")
	}

	return nil
}
