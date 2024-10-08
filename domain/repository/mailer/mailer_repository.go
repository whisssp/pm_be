package mailer

import "mime/multipart"

type MailerRepository interface {
	SendEmailWithPlainText(string, string, []string, multipart.File) error
	SendEmailWithAttachment(string, string, string, []string, string) error
}