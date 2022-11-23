package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/SaidovZohid/auth-redis-verify/config"
)

type SendEmailRequest struct {
	To      []string
	Body    map[string]string
	Subject string
}

func SendEmail(cfg *config.Config, req *SendEmailRequest) error {
	from := cfg.Smtp.Sender
	to := req.To

	password := cfg.Smtp.Password

	var body bytes.Buffer

	t, err := template.ParseFiles("./templates/email_verification.html")
	if err != nil {
		return err
	}
	t.Execute(&body, req.Body)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := fmt.Sprintf("Subject: %s\n", req.Subject)
	msg := []byte(subject + mime + body.String())

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		return err
	}
	return nil
}