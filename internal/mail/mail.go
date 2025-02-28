package mail

import (
    "fmt"
    "net/smtp"

    "github.com/hard-gainer/team-manager/internal/config"
)

type Mailer struct {
    config *config.Config
}

func NewMailer(config *config.Config) *Mailer {
    return &Mailer{
        config: config,
    }
}

func (m *Mailer) SendInvitation(email, inviteLink string) error {
    auth := smtp.PlainAuth("", m.config.SMTPFrom, m.config.SMTPPassword, m.config.SMTPHost)

    subject := "Invitation to join project"
    body := fmt.Sprintf(`
        <html>
        <body>
            <h2>You've been invited to join a project</h2>
            <p>Click the link below to join the project:</p>
            <p><a href="%s">Join Project</a></p>
            <p>This link will expire in 7 days.</p>
        </body>
        </html>
    `, inviteLink)

    msg := fmt.Sprintf("From: %s\r\n"+
        "To: %s\r\n"+
        "MIME-Version: 1.0\r\n"+
        "Content-Type: text/html; charset=UTF-8\r\n"+
        "Subject: %s\r\n\r\n"+
        "%s", m.config.SMTPFrom, email, subject, body)

    return smtp.SendMail(
        fmt.Sprintf("%s:%s", m.config.SMTPHost, m.config.SMTPPort),
        auth,
        m.config.SMTPFrom,
        []string{email},
        []byte(msg),
    )
}