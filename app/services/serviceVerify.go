package services

import (
	"net/smtp"
	"strconv"
)

type EmailService interface {
	SendVerificationEmail(email, verificationLink string) error
}

type smtpEmailService struct {
	smtpServer   string
	smtpPort     int
	smtpUsername string
	smtpPassword string
}

func NewSMTPEmailService(smtpServer string, smtpPort int, smtpUsername, smtpPassword string) EmailService {
	return &smtpEmailService{smtpServer, smtpPort, smtpUsername, smtpPassword}
}

func (s *smtpEmailService) SendVerificationEmail(email, verificationLink string) error {
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpServer)

	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: Verification Email\r\n" +
		"\r\n" +
		"Click the following link to reset your password: " + verificationLink + "\r\n")

	err := smtp.SendMail(s.smtpServer+":"+strconv.Itoa(s.smtpPort), auth, s.smtpUsername, to, msg)
	if err != nil {
		return err
	}

	return nil
}
