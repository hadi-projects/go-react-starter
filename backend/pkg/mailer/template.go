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
<body style="margin:0;padding:0;background-color:#f8fafc;font-family:'Inter', 'SF Pro Text', -apple-system, blinkmacsystemfont, 'Segoe UI', roboto, oxygen, ubuntu, cantarell, 'Open Sans', 'Helvetica Neue', sans-serif; -webkit-font-smoothing: antialiased;">
  <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#f8fafc;padding:48px 0;">
    <tr>
      <td align="center">
        <!-- Main Container -->
        <table width="600" cellpadding="0" cellspacing="0" border="0" style="max-width:600px;width:100%%;background-color:#ffffff;border-radius:24px;box-shadow:0 20px 25px -5px rgba(0,0,0,0.05), 0 10px 10px -5px rgba(0,0,0,0.01);overflow:hidden;">
          
          <!-- Header Area -->
          <tr>
            <td align="center" style="background: linear-gradient(135deg, #4f46e5 0%%, #7c3aed 100%%); padding: 56px 40px 48px;">
                <div style="background: rgba(255, 255, 255, 0.15); width: 64px; height: 64px; border-radius: 20px; display: inline-flex; align-items: center; justify-content: center; margin-bottom: 24px; border: 1px solid rgba(255, 255, 255, 0.2);">
                    <span style="font-size: 32px; line-height: 1;">🔐</span>
                </div>
                <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 800; letter-spacing: -0.025em; line-height: 1.2;">Reset Your Password</h1>
                <p style="margin: 12px 0 0; color: rgba(255, 255, 255, 0.9); font-size: 16px; font-weight: 400;">Secure your account with a new password</p>
            </td>
          </tr>

          <!-- Body Content -->
          <tr>
            <td style="padding: 48px 48px 40px;">
              <p style="margin: 0 0 24px; color: #1e293b; font-size: 18px; font-weight: 600;">Hello there,</p>
              <p style="margin: 0 0 32px; color: #475569; font-size: 16px; line-height: 1.7;">
                We received a request to reset the password for your account. If you didn't make this request, you can safely ignore this email. Otherwise, click the button below to choose a new password:
              </p>

              <!-- Action Button -->
              <table width="100%%" border="0" cellspacing="0" cellpadding="0">
                <tr>
                  <td align="center">
                    <a href="%s" style="display: inline-block; background: #4f46e5; color: #ffffff; font-size: 16px; font-weight: 700; text-decoration: none; padding: 18px 44px; border-radius: 14px; box-shadow: 0 10px 15px -3px rgba(79, 70, 229, 0.3); transition: all 0.2s ease;">
                      Reset Password Now
                    </a>
                  </td>
                </tr>
              </table>

              <p style="margin: 40px 0 0; color: #64748b; font-size: 14px; text-align: center;">
                This link will expire in <strong style="color: #1e293b;">15 minutes</strong> for your security.
              </p>

              <!-- Divider -->
              <div style="height: 1px; background-color: #f1f5f9; margin: 40px 0 32px;"></div>

              <!-- Manual Link -->
              <p style="margin: 0 0 12px; color: #94a3b8; font-size: 12px; font-weight: 500; text-transform: uppercase; letter-spacing: 0.05em;">If the button doesn't work, copy this link:</p>
              <div style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 12px 16px; word-break: break-all;">
                <a href="%s" style="color: #4f46e5; font-size: 13px; text-decoration: none; font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;">%s</a>
              </div>
            </td>
          </tr>

          <!-- Footer Area -->
          <tr>
            <td style="background-color: #f8fafc; padding: 32px 48px; border-top: 1px solid #edf2f7; text-align: center;">
              <p style="margin: 0; color: #94a3b8; font-size: 13px; line-height: 1.5;">
                Sent with ❤️ from <strong>Go-React Starter Team</strong>.<br>
                &copy; 2026 Admin Panel Inc. All rights reserved.
              </p>
            </td>
          </tr>
        </table>

        <!-- Security Footer -->
        <table width="600" cellpadding="0" cellspacing="0" border="0" style="max-width:600px;width:100%%;">
            <tr>
                <td style="padding-top: 24px;">
                    <p style="margin: 0; color: #94a3b8; font-size: 11px; text-align: center; line-height: 1.6;">
                        You received this because a password reset was requested for your account.<br>
                        If you believe this was an error, please contact security@yourdomain.com
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

func GetTwoFAResetEmailNative(resetLink string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Reset 2FA/OTP</title>
</head>
<body style="margin:0;padding:0;background-color:#f8fafc;font-family:'Inter', 'SF Pro Text', -apple-system, blinkmacsystemfont, 'Segoe UI', roboto, oxygen, ubuntu, cantarell, 'Open Sans', 'Helvetica Neue', sans-serif; -webkit-font-smoothing: antialiased;">
  <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#f8fafc;padding:48px 0;">
    <tr>
      <td align="center">
        <table width="600" cellpadding="0" cellspacing="0" border="0" style="max-width:600px;width:100%%;background-color:#ffffff;border-radius:24px;box-shadow:0 20px 25px -5px rgba(0,0,0,0.05), 0 10px 10px -5px rgba(0,0,0,0.01);overflow:hidden;">
          <tr>
            <td align="center" style="background: linear-gradient(135deg, #10b981 0%%, #059669 100%%); padding: 56px 40px 48px;">
                <div style="background: rgba(255, 255, 255, 0.15); width: 64px; height: 64px; border-radius: 20px; display: inline-flex; align-items: center; justify-content: center; margin-bottom: 24px; border: 1px solid rgba(255, 255, 255, 0.2);">
                    <span style="font-size: 32px; line-height: 1;">🛡️</span>
                </div>
                <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 800; letter-spacing: -0.025em; line-height: 1.2;">Reset Your 2FA</h1>
                <p style="margin: 12px 0 0; color: rgba(255, 255, 255, 0.9); font-size: 16px; font-weight: 400;">Disable 2FA to regain access to your account</p>
            </td>
          </tr>
          <tr>
            <td style="padding: 48px 48px 40px;">
              <p style="margin: 0 0 24px; color: #1e293b; font-size: 18px; font-weight: 600;">Hello there,</p>
              <p style="margin: 0 0 32px; color: #475569; font-size: 16px; line-height: 1.7;">
                We received a request to bypass and disable Two-Factor Authentication (2FA) for your account because you lost access to your authenticator app. Click the button below to confirm and disable 2FA:
              </p>
              <table width="100%%" border="0" cellspacing="0" cellpadding="0">
                <tr>
                  <td align="center">
                    <a href="%s" style="display: inline-block; background: #10b981; color: #ffffff; font-size: 16px; font-weight: 700; text-decoration: none; padding: 18px 44px; border-radius: 14px; box-shadow: 0 10px 15px -3px rgba(16, 185, 129, 0.3); transition: all 0.2s ease;">
                      Disable 2FA Now
                    </a>
                  </td>
                </tr>
              </table>
              <p style="margin: 40px 0 0; color: #64748b; font-size: 14px; text-align: center;">
                This link will expire in <strong style="color: #1e293b;">15 minutes</strong> for your security.
              </p>
              <div style="height: 1px; background-color: #f1f5f9; margin: 40px 0 32px;"></div>
              <p style="margin: 0 0 12px; color: #94a3b8; font-size: 12px; font-weight: 500; text-transform: uppercase; letter-spacing: 0.05em;">If the button doesn't work, copy this link:</p>
              <div style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 12px 16px; word-break: break-all;">
                <a href="%s" style="color: #10b981; font-size: 13px; text-decoration: none; font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;">%s</a>
              </div>
            </td>
          </tr>
          <tr>
            <td style="background-color: #f8fafc; padding: 32px 48px; border-top: 1px solid #edf2f7; text-align: center;">
              <p style="margin: 0; color: #94a3b8; font-size: 13px; line-height: 1.5;">
                Sent with ❤️ from <strong>Go-React Starter Team</strong>.<br>
                &copy; 2026 Admin Panel Inc. All rights reserved.
              </p>
            </td>
          </tr>
        </table>
        <table width="600" cellpadding="0" cellspacing="0" border="0" style="max-width:600px;width:100%%;">
            <tr>
                <td style="padding-top: 24px;">
                    <p style="margin: 0; color: #94a3b8; font-size: 11px; text-align: center; line-height: 1.6;">
                        You received this because a 2FA reset was requested for your account.<br>
                        If you believe this was an error, please contact security@yourdomain.com
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
