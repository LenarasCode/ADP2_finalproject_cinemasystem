package email

import (
	"net/smtp"
	"os"
)

type Sender struct {
	from string
	auth smtp.Auth
	host string
	port string
}

func NewSender() *Sender {
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	host := "smtp.gmail.com"
	port := "587"
	return &Sender{
		from: from,
		auth: smtp.PlainAuth("", from, pass, host),
		host: host,
		port: port,
	}
}

func (s *Sender) Send(to, subject, body string) error {
	msg := []byte("To: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + body)
	return smtp.SendMail(s.host+":"+s.port, s.auth, s.from, []string{to}, msg)
}
