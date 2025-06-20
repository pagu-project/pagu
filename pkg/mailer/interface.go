package mailer

type IMailer interface {
	SendTemplateMailAsync(recipient, templatePath string, data map[string]string) error
	SendTemplateMail(recipient, templatePath string, data map[string]string) error
}
