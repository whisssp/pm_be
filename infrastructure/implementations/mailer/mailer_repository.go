package mailer

import (
	"fmt"
	gomail "gopkg.in/mail.v2"
	"mime/multipart"
	"pm/domain/repository/mailer"
	"pm/infrastructure/persistences/base"
)

type MailerRepository struct {
	p *base.Persistence
}

func (m MailerRepository) SendEmailWithPlainText(body, subject string, receivers []string, file multipart.File) error {
	mM := gomail.NewMessage()

	mM.SetHeader("To", receivers...)
	mM.SetHeader("From", m.p.Mailer.From)
	mM.SetHeader("Subject", subject)
	mM.SetBody("text/plain", body)
	if err := m.p.Mailer.D.DialAndSend(mM); err != nil {
		fmt.Println("Error send mail: " + err.Error())
		return err
	}

	return nil
}

func NewMailerRepository(p *base.Persistence) mailer.MailerRepository {
	return &MailerRepository{p: p}
}