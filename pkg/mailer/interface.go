package mailer

type IMailer interface {
	SendTemplateMailAsync(email, templatePath string, data map[string]string) error
	SendTemplateMail(email, templatePath string, data map[string]string) error
}
