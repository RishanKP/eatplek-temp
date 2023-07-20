package services

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendPasswordResetMail(email, name string, token string) error {
	from := mail.NewEmail("Eatplek", os.Getenv("FROM_EMAIL"))
	subject := "Password Reset"
	to := mail.NewEmail(name, email)
	htmlContent := "<strong>Hello, " + name + "! Click on the link below to reset your password: https://restaurant.eatplek.com/reset?token=" + token + "</strong><br> Link Expires in 15 minutes. "

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}

func SendVerificationEmail(email, name string, token string) error {
	from := mail.NewEmail("Eatplek", os.Getenv("FROM_EMAIL"))
	subject := "Verify Account"
	to := mail.NewEmail(name, email)
	htmlContent := "<strong>Hello, " + name + "! Click on the link below to verify your Eatplek account: https://user.eatplek.com/email-verification?token=" + token + "</strong><br> Link Expires in 15 minutes. "

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
