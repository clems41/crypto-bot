package gmail

import (
	"crypto-bot/external/service/mailService"
	"crypto-bot/pkg/utils/envUtils"
	"net/smtp"
)

var _ mailService.Service = (*service)(nil)

type service struct {
	smtpUser     string
	smtpPassword string
}

func NewService() (svc *service, err error) {
	smtpUser, err := envUtils.GetFromEnvOrError(envSmtpUser)
	if err != nil {
		return
	}
	smtpPassword, err := envUtils.GetFromEnvOrError(envSmtpPassword)
	if err != nil {
		return
	}

	svc = &service{
		smtpUser:     smtpUser,
		smtpPassword: smtpPassword,
	}

	return
}

func (svc *service) Send(request mailService.SendRequest) (response mailService.SendResponse, err error) {
	// Authentication.
	auth := smtp.PlainAuth("", svc.smtpUser, svc.smtpPassword, smtpServer)

	// Sending email.
	err = smtp.SendMail(smtpServer+":"+smtpPortTls, auth, request.From, []string{svc.smtpUser}, []byte(request.Message))
	if err != nil {
		return
	}
	return
}
