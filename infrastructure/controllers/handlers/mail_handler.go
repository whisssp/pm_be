package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pm/application"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
)

type MailHandler struct {
	p             *base.Persistence
	mailerUseCase application.MailUsecase
}

func NewMailHandler(p *base.Persistence) *MailHandler {
	mailerUsecase := application.NewMailUsecase(p)
	return &MailHandler{p: p, mailerUseCase: mailerUsecase}
}

func (h *MailHandler) HandleSendEmail(c *gin.Context) {
	var mailRequest payload.MailRequest
	if err := c.ShouldBindJSON(&mailRequest); err != nil {
		c.Error(err)
	}

	message, err := h.mailerUseCase.SendMail(c, &mailRequest)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, payload.SuccessResponse(nil, message))
}