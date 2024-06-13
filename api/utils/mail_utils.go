package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strings"
)

func SendEmailHTML(subject string, to []string, path string, params any) error {
	var body bytes.Buffer
	var err error

	t, err := template.ParseFiles(path)

	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("ocurrio un error inesperado")
	}

	t.Execute(&body, params)

	auth := smtp.PlainAuth(
		os.Getenv("MAIL_USERNAME"),
		os.Getenv("MAIL_USERNAME"),
		os.Getenv("MAIL_PASSWORD"),
		os.Getenv("MAIL_HOST"),
	)

	msg := []byte("To: " + to[0:1][0] + "\r\n" +
		"Cc: " + strings.Join(to[1:], ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";" + "\n\n" +
		body.String())

	if err = smtp.SendMail(
		os.Getenv("MAIL_HOST")+":"+os.Getenv("MAIL_PORT"),
		auth,
		os.Getenv("MAIL_FROM_ADDRESS"),
		to,
		[]byte(msg),
	); err != nil {
		log.Println(err.Error())
		return fmt.Errorf("ocurrio un error inesperado")
	}

	return nil
}
