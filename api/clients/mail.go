package clients

import (
	"bytes"
	"context"
	"fmt"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/mailersend/mailersend-go"
	"html/template"
	"os"
	"time"
)

//func SendMail(templateName, to, subject string, data interface{}) {
//	t, err := template.ParseFiles(fmt.Sprintf("templates/%v", templateName))
//	if err != nil {
//		utils.Logger.Error("error parsing html file", err)
//		return
//	}
//
//	var buff bytes.Buffer
//
//	err = t.Execute(&buff, data)
//	if err != nil {
//		utils.Logger.Error("error executing html file", err)
//		return
//	}
//
//	m := gomail.NewMessage()
//	m.SetHeader("From", "no-reply@daryanamin.co.uk")
//	m.SetHeader("To", to)
//	m.SetHeader("Subject", subject)
//	m.SetBody("text/html", buff.String())
//
//	portStr := os.Getenv("SMTP_PORT")
//
//	port, _ := strconv.ParseInt(portStr, 10, 64)
//
//	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), int(port), os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
//	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
//
//	if err = d.DialAndSend(m); err != nil {
//		utils.Logger.Error("error sending email", err)
//	}
//}

func SendMail(templateName, to, subject string, data interface{}) {
	t, err := template.ParseFiles(fmt.Sprintf("templates/%v", templateName))
	if err != nil {
		utils.Logger.Error("error parsing html file", err)
		return
	}

	var buff bytes.Buffer

	err = t.Execute(&buff, data)
	if err != nil {
		utils.Logger.Error("error executing html file", err)
		return
	}

	APIKey := os.Getenv("MAILERSEND_API_KEY")
	ms := mailersend.NewMailersend(APIKey)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	html := buff.String()

	from := mailersend.From{
		Name:  "Daryan Amin",
		Email: "no-reply@daryanamin.co.uk",
	}

	recipients := []mailersend.Recipient{
		{
			Name:  to,
			Email: to,
		},
	}

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)

	if _, err = ms.Email.Send(ctx, message); err != nil {
		utils.Logger.Error("error sending email", err)
	}
}
