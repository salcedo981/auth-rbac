package hlpAuth

import (
	"fmt"
	"net/smtp"

	utils_v1 "github.com/FDSAP-Git-Org/hephaestus/utils/v1"
)

func SendTempPasswordEmail(toEmail, username, tempPassword string) error {
	from := utils_v1.GetEnv("SMTP_USER")
	smtpHost := utils_v1.GetEnv("SMTP_HOST")
	password := utils_v1.GetEnv("SMTP_PASS")
	smtpPort := utils_v1.GetEnv("SMTP_PORT")

	subject := "Subject: Your Temporary Account Credentials\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin:0;padding:0;background-color:#f4f6f8;font-family:Arial,Helvetica,sans-serif;">
  <div style="max-width:600px;margin:0 auto;background:#ffffff;border-radius:8px;overflow:hidden;">
    
    <div style="background:#1f2937;color:#ffffff;padding:20px;text-align:center;">
      <h1 style="margin:0;font-size:22px;">Account Created</h1>
    </div>

    <div style="padding:20px;color:#111827;font-size:14px;line-height:1.6;">
      <p>Hi <strong>%s</strong>,</p>

      <p>Your account has been successfully created. Below are your temporary login credentials:</p>

      <div style="background:#f3f4f6;padding:16px;border-radius:6px;margin:15px 0;">
        <p style="margin:0;"><strong>Username:</strong> %s</p>
        <p style="margin:0;"><strong>Temporary Password:</strong> %s</p>
      </div>

      <p style="color:#b91c1c;">
        For security reasons, please change your password immediately after logging in.
      </p>

      <p style="margin-top:20px;">Thank you,<br><strong>Security Team</strong></p>
    </div>

    <div style="background:#f9fafb;text-align:center;padding:15px;font-size:12px;color:#6b7280;">
      <p style="margin:0;">© 2025 Your Company. All rights reserved.</p>
    </div>

  </div>
</body>
</html>
`, username, username, tempPassword)

	message := []byte(subject + mime + htmlBody)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{toEmail},
		message,
	)
}
