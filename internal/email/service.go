package email

import (
	"fmt"
	"net/smtp"

	"github.com/Jeffreasy/GoBackend/configs"
)

type Service interface {
	SendMail(to, subject, body string) error
	FromEmail() string
}

type emailService struct {
	cfg *configs.Config
}

func NewService(cfg *configs.Config) Service {
	return &emailService{cfg: cfg}
}

func (e *emailService) SendMail(to, subject, body string) error {
	auth := smtp.PlainAuth("", e.cfg.SMTPUser, e.cfg.SMTPPassword, e.cfg.SMTPHost)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))
	addr := fmt.Sprintf("%s:%s", e.cfg.SMTPHost, e.cfg.SMTPPort)
	return smtp.SendMail(addr, auth, e.cfg.FromEmail, []string{to}, msg)
}

func (e *emailService) FromEmail() string {
	return e.cfg.FromEmail
}
