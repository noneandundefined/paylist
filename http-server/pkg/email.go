package pkg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
	"paylist.server/infra/locale"
	"paylist.server/infra/logger"
)

/* Отправка email письма на почту */
func SendEmail(to, title, content string, tr locale.Translator) error {
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
	mail.SetBody("text/html", BuildEmailTemplate(content, tr))

	d := gomail.NewDialer(smtpAddr, port, smtpEmail, smtpPassword)
	d.SSL = port == 465

	if err := d.DialAndSend(mail); err != nil {
		return fmt.Errorf("%s: %v", tr.TErr("error.email-send-failed"), err)
	}

	logger.Info("SendEmail email={%s}: Email has been sent!", to)
	return nil
}

func BuildEmailTemplate(html string, tr locale.Translator) string {
	clientURL := strings.TrimRight(strings.TrimSpace(os.Getenv("CLIENT_URL")), "/")
	logoURL := clientURL + "/local/images/favicon-512.jpg"

	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
		</head>
		<body style="margin:0; padding:0; font-family:Arial,Helvetica,sans-serif; font-size:16px; background-color:#ffffff;">
			<div style="max-width:625px; margin:40px auto; background-color:#ffffff; overflow:hidden;">
				<div style="text-align:center;margin-bottom:40px;">
					<img src="%s" alt="Paylist" width="64" style="display:block;margin:0 auto 12px;" />
				</div>

				<div style="color:#222222; line-height:200%%;">
					%s

					<p>%s</p>
				</div>

				<hr style="border:none;border-top:1px solid #e9edf2;margin:30px 0;">

				<div style="text-align:center; font-size:13px; color:#222222;">
					<img src="%s" alt="Paylist" width="40" style="display:block;margin:0 auto;" />

					<p>%s</p>

					<p style="color:#adb1b8;">%s</p>
					<p style="color:#adb1b8;">%s</p>
				</div>
			</div>
		</body>
		</html>
	`,
		logoURL,
		html,
		tr.T("email-regards"),
		logoURL,
		tr.T("email-stay-informed"),
		tr.T("email-created-automatic"),
		fmt.Sprintf(tr.T("copyright"), time.Now().Year()),
	)
}
