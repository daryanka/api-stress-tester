package clients

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/daryanka/api-stress-tester/api/utils"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"strconv"
)

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

	m := gomail.NewMessage()
	m.SetHeader("From", "daryanka@hotmail.co.uk")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", buff.String())

	portStr := os.Getenv("SMTP_PORT")

	port, _ := strconv.ParseInt(portStr, 10, 64)

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), int(port), os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err = d.DialAndSend(m); err != nil {
		utils.Logger.Error("error sending email", err)
	}
}