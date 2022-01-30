package gmail

import (
	"crypto-bot/external/service/mailService"
	"crypto-bot/pkg/logger"
	"crypto-bot/pkg/utils/envUtils"
	"fmt"
	"net/smtp"
)

var _ mailService.Service = (*service)(nil)

type service struct {
	smtpUser     string
	smtpPassword string
	sendMail     bool
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
	sendMailStr := envUtils.GetFromEnvOrDefault(envSendMail, defaultSendMail)

	svc = &service{
		smtpUser:     smtpUser,
		smtpPassword: smtpPassword,
		sendMail:     sendMailStr == "true",
	}

	return
}

func (svc *service) Send(request mailService.SendRequest) (response mailService.SendResponse, err error) {
	if !svc.sendMail {
		logger.Debugf("SEND_MAIL variable has been set to false, mail will not be sent.")
		return
	}
	// Authentication.
	auth := smtp.PlainAuth("", svc.smtpUser, svc.smtpPassword, smtpServer)
	message := fmt.Sprintf("From: <%s>\r\nTo: <%s>\r\nSubject: %s\r\n\r\n%s",
		request.From, svc.smtpUser, request.Subject, request.Body)

	// Sending email.
	err = smtp.SendMail(smtpServer+":"+smtpPortTls, auth, request.From, []string{svc.smtpUser}, []byte(message))
	if err != nil {
		return
	}
	return
}
