package mailer

import "fmt"

func GetResetPasswordEmailNative(resetLink string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Reset Your Password</title>
</head>
<body style="margin:0;padding:0;background-color:#f4f6f9;font-family:'Segoe UI',Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f6f9;padding:40px 0;">
    <tr>
      <td align="center">
        <table width="600" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%%;">

          <!-- Header -->
          <tr>
            <td align="center" style="background:linear-gradient(135deg,#667eea 0%%,#764ba2 100%%);border-radius:12px 12px 0 0;padding:40px 40px 32px;">
              <div style="width:56px;height:56px;background:rgba(255,255,255,0.2);border-radius:50%%;display:inline-flex;align-items:center;justify-content:center;margin-bottom:16px;">
                <span style="font-size:28px;">🔐</span>
              </div>
              <h1 style="margin:0;color:#ffffff;font-size:26px;font-weight:700;letter-spacing:-0.5px;">Reset Your Password</h1>
              <p style="margin:8px 0 0;color:rgba(255,255,255,0.8);font-size:14px;">We received a request to reset your password</p>
            </td>
          </tr>

          <!-- Body -->
          <tr>
            <td style="background:#ffffff;padding:40px;">
              <p style="margin:0 0 20px;color:#374151;font-size:16px;line-height:1.6;">Hi there,</p>
              <p style="margin:0 0 28px;color:#6b7280;font-size:15px;line-height:1.7;">
                Someone requested a password reset for your account. If this was you, click the button below to set a new password. This link is valid for <strong style="color:#374151;">15 minutes</strong>.
              </p>

              <!-- CTA Button -->
              <table width="100%%" cellpadding="0" cellspacing="0">
                <tr>
                  <td align="center" style="padding:8px 0 32px;">
                    <a href="%s" style="display:inline-block;background:linear-gradient(135deg,#667eea 0%%,#764ba2 100%%);color:#ffffff;text-decoration:none;font-size:16px;font-weight:600;padding:14px 40px;border-radius:8px;letter-spacing:0.3px;box-shadow:0 4px 15px rgba(102,126,234,0.4);">
                      Reset Password →
                    </a>
                  </td>
                </tr>
              </table>

              <!-- Fallback link -->
              <p style="margin:0 0 8px;color:#9ca3af;font-size:13px;">If the button doesn't work, copy and paste this link into your browser:</p>
              <p style="margin:0 0 28px;word-break:break-all;">
                <a href="%s" style="color:#667eea;font-size:13px;text-decoration:none;">%s</a>
              </p>

              <!-- Divider -->
              <hr style="border:none;border-top:1px solid #e5e7eb;margin:0 0 24px;">

              <!-- Security notice -->
              <table width="100%%" cellpadding="0" cellspacing="0" style="background:#fef9c3;border:1px solid #fde68a;border-radius:8px;">
                <tr>
                  <td style="padding:16px 20px;">
                    <p style="margin:0;color:#92400e;font-size:13px;line-height:1.6;">
                      ⚠️ <strong>Didn't request this?</strong> You can safely ignore this email. Your password will remain unchanged and no action is required.
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="background:#f9fafb;border-radius:0 0 12px 12px;padding:24px 40px;border-top:1px solid #e5e7eb;">
              <p style="margin:0;color:#9ca3af;font-size:12px;text-align:center;line-height:1.6;">
                This email was sent by <strong style="color:#6b7280;">Go-React Starter</strong>.<br>
                For security, this link expires in 15 minutes and can only be used once.
              </p>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, resetLink, resetLink, resetLink)
}
