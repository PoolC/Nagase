package email

import (
	"fmt"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var client *sendgrid.Client

var emailHTMLTemplate = `
<html>
  <head>
    <meta charset="utf-8">
  </head>
  <body style="margin: 0; font-family: 'Arial', 'san-serif';">
    <div style="height: 180px; background: #009688; color: #ffffff; text-align: center;">
      <h1 style="font-size: 360%%; font-weight: 100; padding-top: 35px; margin-top: 0; margin-bottom: 10px;">PoolC</h1>
      <h2 style="font-size: 180%%; font-weight: 200; margin: 0px;">%s</h2>
    </div>
    <div style="padding: 20px;">
      <p>%s</p>
    </div>
  </body>
</html>
`

type Email struct {
	Title string
	Body  string
	To    string
}

func (e *Email) Send() error {
	sgMail := mail.NewSingleEmail(mail.NewEmail("PoolC", "noreply@poolc.org"), "[PoolC] " + e.Title, mail.NewEmail("PoolC", e.To), e.Body, e.bodyHTML())
	_, err := client.Send(sgMail)
	return err
}

func (e *Email) bodyHTML() string {
	return fmt.Sprintf(emailHTMLTemplate, e.Title, strings.Replace(e.Body, "\n", "<br/>", -1))
}

func init() {
	client = sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
}
