package mail

import (
	"crypto/tls"
	"github.com/joho/godotenv"
	"gopkg.in/mail.v2"
	"pm/infrastructure/config"
)

type Mailer struct {
	D    *mail.Dialer
	From string
}

func NewMailer() *Mailer {
	godotenv.Load()
	username := config.GetEnv("pm.mail.username", "nghia14802@gmail.com")
	password := config.GetEnv("pm.mail.password", "")
	port := config.GetEnvAsInt("pm.mail.smtp.port", 587)
	host := config.GetEnv("pm.mail.smtp.host", "smtp.gmail.com")

	d := mail.NewDialer(host, int(port), username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &Mailer{D: d, From: username}
}