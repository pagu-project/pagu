package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
	"github.com/pagu-project/pagu/pkg/log"
)

type SMTPMailer struct {
	dialer *mail.Dialer
	cfg    *Config
}

func NewSMTPMailer(cfg *Config) *SMTPMailer {
	dialer := mail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	dialer.Timeout = 5 * time.Second

	return &SMTPMailer{
		dialer: dialer,
		cfg:    cfg,
	}
}

func (*SMTPMailer) LoadMailTemplate(path string) (*template.Template, error) {
	tmpl, err := template.New("").ParseFiles(path)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (s *SMTPMailer) SendTemplateMailAsync(recipient, tmplPath string, data map[string]string) error {
	msg, err := s.makeMessage(recipient, tmplPath, data)
	if err != nil {
		return err
	}

	go func() {
		err = s.dialer.DialAndSend(msg)
		if err != nil {
			log.Warn("failed to send voucher email", "error", err)
		}
	}()

	return nil
}

func (s *SMTPMailer) SendTemplateMail(recipient, tmplPath string, data map[string]string) error {
	msg, err := s.makeMessage(recipient, tmplPath, data)
	if err != nil {
		return err
	}

	return s.dialer.DialAndSend(msg)
}

func (s *SMTPMailer) makeMessage(recipient, tmplPath string, data map[string]string) (*mail.Message, error) {
	tmpl, err := s.LoadMailTemplate(tmplPath)
	if err != nil {
		return nil, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return nil, fmt.Errorf("error executing template with subject: %w", err)
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return nil, fmt.Errorf("error executing plain body: %w", err)
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return nil, fmt.Errorf("error executing HTML body: %w", err)
	}

	msg := mail.NewMessage()
	msg.SetHeader("From", s.cfg.Sender)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	return msg, nil
}
