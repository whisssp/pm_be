package application

import (
	"github.com/gin-gonic/gin"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/implementations/mailer"
	"pm/infrastructure/persistences/base"
)

type MailUsecase interface {
	SendMail(*gin.Context, *payload.MailRequest) (string, error)
}

type mailUsecase struct {
	p *base.Persistence
}

func (m mailUsecase) SendMail(c *gin.Context, mailRequest *payload.MailRequest) (string, error) {

	mailerRepo := mailer.NewMailerRepository(m.p)
	err := mailerRepo.SendEmailWithPlainText(mailRequest.Body, mailRequest.Subject, mailRequest.Receivers, nil)
	if err != nil {
		return "", err
	}
	return "dang gui doi xiu", nil
}

func NewMailUsecase(p *base.Persistence) MailUsecase {
	return &mailUsecase{p: p}
}