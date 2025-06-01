package mailer

type IMailer interface {
	SendTemplateMail(recipient, templatePath string, data map[string]string) error
}
