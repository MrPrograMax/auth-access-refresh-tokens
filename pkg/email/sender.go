package email

import "fmt"

type SendEmailInput struct {
	To      string
	Subject string
	Body    string
}

type Sender interface {
	Send(input SendEmailInput) error
}

func GenerateIpWarningEmail(to string, newIP string) SendEmailInput {
	subject := "Замечена подозрительная активность"
	body := fmt.Sprintf("Мы заметили вход с нового IP адреса: %s.\nЕсли это были не вы обратитесь в тех поддержку.", newIP)

	return SendEmailInput{
		To:      to,
		Subject: subject,
		Body:    body,
	}
}
