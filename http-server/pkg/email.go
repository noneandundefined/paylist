package pkg

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
	"paylist.server/infra/locale"
	"paylist.server/infra/logger"
)

/* Отправка email письма на почту */
func SendEmail(to, title, content string, tr locale.Translator) error {
	go_env := os.Getenv("GO_ENV") == "DEV"

	port, err := strconv.Atoi(strings.TrimSpace(os.Getenv("SMTP_PORT")))
	if err != nil || port == 0 {
		port = 587
	}

	smtpAddr := strings.Trim(os.Getenv("SMTP_ADDR"), `"' `)
	smtpEmail := strings.Trim(os.Getenv("SMTP_EMAIL"), `"' `)
	smtpPassword := strings.Trim(os.Getenv("SMTP_PASSWORD"), `"' `)

	mail := gomail.NewMessage()
	mail.SetHeader("From", smtpEmail)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", title)
	mail.SetBody("text/html", content)

	d := gomail.NewDialer(smtpAddr, port, smtpEmail, smtpPassword)
	// Port 465 — implicit SSL; port 587 — STARTTLS (SSL must be false).
	d.SSL = port == 465

	if go_env {
		logger.Info("SendEmail: SMTP %s:%d ssl=%v", smtpAddr, port, d.SSL)
	}

	if err := d.DialAndSend(mail); err != nil {
		return fmt.Errorf("%s: %v", tr.TErr("error.email-send-failed"), err)
	}

	logger.Info("SendEmail email={%s}: Email has been sent!", to)
	return nil
}
