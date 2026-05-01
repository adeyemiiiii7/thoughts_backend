package email

import (
	"fmt"
	"os"
	"strconv"

	mail "github.com/wneessen/go-mail"
)

func SendVerificationEmail(to, token string) error {
	host := os.Getenv("SMTP_HOST")
	portValue := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	fromName := os.Getenv("SMTP_FROM_NAME")
	backendURL := os.Getenv("BACKEND_URL")

	if host == "" || portValue == "" || username == "" || password == "" || fromEmail == "" {
		return fmt.Errorf("SMTP settings are incomplete")
	}

	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}

	port, err := strconv.Atoi(portValue)
	if err != nil {
		return fmt.Errorf("invalid SMTP_PORT: %w", err)
	}

	client, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	message := mail.NewMsg()
	from := fromEmail
	if fromName != "" {
		from = fmt.Sprintf("%s <%s>", fromName, fromEmail)
	}

	if err := message.From(from); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}

	if err := message.To(to); err != nil {
		return fmt.Errorf("failed to set recipient address: %w", err)
	}

	verifyURL := fmt.Sprintf("%s/auth/verify-email?token=%s", backendURL, token)

	message.Subject("Verify your email address")
	message.SetBodyString(mail.TypeTextPlain, fmt.Sprintf(
		"Verify your email for Thoughts.\n\nOpen this link to verify your account:\n%s\n\nThis link expires in 24 hours.\nIf you did not create an account, you can ignore this email.\n",
		verifyURL,
	))
	message.AddAlternativeString(mail.TypeTextHTML, fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
  <body style="margin:0; padding:0; background-color:#f5f5f5; color:#111111; font-family:Arial, Helvetica, sans-serif;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background-color:#f5f5f5; padding:32px 16px;">
      <tr>
        <td align="center">
          <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:560px; background-color:#ffffff; border:1px solid #e5e5e5;">
            <tr>
              <td style="background-color:#111111; padding:20px 24px;">
                <p style="margin:0; color:#ffffff; font-size:14px; letter-spacing:1px; text-transform:uppercase;">
                  Thoughts
                </p>
              </td>
            </tr>
            <tr>
              <td style="padding:40px 24px 24px 24px;">
                <h1 style="margin:0 0 16px 0; font-size:28px; line-height:1.2; color:#111111;">
                  Verify your email
                </h1>
                <p style="margin:0 0 16px 0; font-size:16px; line-height:1.6; color:#333333;">
                  Thanks for signing up. Confirm your email address to activate your account and continue.
                </p>
                <p style="margin:0 0 32px 0; font-size:16px; line-height:1.6; color:#333333;">
                  This verification link will expire in 24 hours.
                </p>
                <table role="presentation" cellspacing="0" cellpadding="0">
                  <tr>
                    <td style="background-color:#111111; border-radius:0;">
                      <a href="%s" style="display:inline-block; padding:14px 24px; font-size:14px; font-weight:bold; color:#ffffff; text-decoration:none; letter-spacing:0.3px;">
                        Verify Email
                      </a>
                    </td>
                  </tr>
                </table>
                <p style="margin:32px 0 12px 0; font-size:14px; line-height:1.6; color:#555555;">
                  If the button does not work, copy and paste this link into your browser:
                </p>
                <p style="margin:0; font-size:13px; line-height:1.6; word-break:break-all; color:#111111;">
                  <a href="%s" style="color:#111111; text-decoration:underline;">
                    %s
                  </a>
                </p>
              </td>
            </tr>
            <tr>
              <td style="padding:24px; border-top:1px solid #e5e5e5;">
                <p style="margin:0 0 8px 0; font-size:13px; line-height:1.6; color:#666666;">
                  If you did not create an account, you can ignore this email.
                </p>
                <p style="margin:0; font-size:12px; line-height:1.6; color:#999999;">
                  Thoughts
                </p>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`, verifyURL, verifyURL, verifyURL))

	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func SendPasswordResetEmail(to, token string) error {
	host := os.Getenv("SMTP_HOST")
	portValue := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	fromName := os.Getenv("SMTP_FROM_NAME")
	backendURL := os.Getenv("BACKEND_URL")

	if host == "" || portValue == "" || username == "" || password == "" || fromEmail == "" {
		return fmt.Errorf("SMTP settings are incomplete")
	}

	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}

	port, err := strconv.Atoi(portValue)
	if err != nil {
		return fmt.Errorf("invalid SMTP_PORT: %w", err)
	}

	client, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	message := mail.NewMsg()
	from := fromEmail
	if fromName != "" {
		from = fmt.Sprintf("%s <%s>", fromName, fromEmail)
	}

	if err := message.From(from); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}

	if err := message.To(to); err != nil {
		return fmt.Errorf("failed to set recipient address: %w", err)
	}

	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", backendURL, token)

	message.Subject("Reset your password")
	message.SetBodyString(mail.TypeTextPlain, fmt.Sprintf(
		"Reset your password for Thoughts.\n\nOpen this link to reset your password:\n%s\n\nThis link expires in 1 hour.\nIf you did not request this, you can ignore this email.\n",
		resetURL,
	))
	message.AddAlternativeString(mail.TypeTextHTML, fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
  <body style="margin:0; padding:0; background-color:#f5f5f5; color:#111111; font-family:Arial, Helvetica, sans-serif;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background-color:#f5f5f5; padding:32px 16px;">
      <tr>
        <td align="center">
          <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:560px; background-color:#ffffff; border:1px solid #e5e5e5;">
            <tr>
              <td style="background-color:#111111; padding:20px 24px;">
                <p style="margin:0; color:#ffffff; font-size:14px; letter-spacing:1px; text-transform:uppercase;">
                  Thoughts
                </p>
              </td>
            </tr>
            <tr>
              <td style="padding:40px 24px 24px 24px;">
                <h1 style="margin:0 0 16px 0; font-size:28px; line-height:1.2; color:#111111;">
                  Reset your password
                </h1>
                <p style="margin:0 0 16px 0; font-size:16px; line-height:1.6; color:#333333;">
                  We received a request to reset your password. Use the button below to continue.
                </p>
                <p style="margin:0 0 32px 0; font-size:16px; line-height:1.6; color:#333333;">
                  This reset link will expire in 1 hour.
                </p>
                <table role="presentation" cellspacing="0" cellpadding="0">
                  <tr>
                    <td style="background-color:#111111; border-radius:0;">
                      <a href="%s" style="display:inline-block; padding:14px 24px; font-size:14px; font-weight:bold; color:#ffffff; text-decoration:none; letter-spacing:0.3px;">
                        Reset Password
                      </a>
                    </td>
                  </tr>
                </table>
                <p style="margin:32px 0 12px 0; font-size:14px; line-height:1.6; color:#555555;">
                  If the button does not work, copy and paste this link into your browser:
                </p>
                <p style="margin:0; font-size:13px; line-height:1.6; word-break:break-all; color:#111111;">
                  <a href="%s" style="color:#111111; text-decoration:underline;">
                    %s
                  </a>
                </p>
              </td>
            </tr>
            <tr>
              <td style="padding:24px; border-top:1px solid #e5e5e5;">
                <p style="margin:0 0 8px 0; font-size:13px; line-height:1.6; color:#666666;">
                  If you did not request a password reset, you can ignore this email.
                </p>
                <p style="margin:0; font-size:12px; line-height:1.6; color:#999999;">
                  Thoughts
                </p>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`, resetURL, resetURL, resetURL))

	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}
