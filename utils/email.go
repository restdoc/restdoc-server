package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/levigross/grequests"

	"restdoc/config"
	"restdoc/render"
)

type html struct {
	Body body `xml:"body"`
}
type body struct {
	Content string `xml:",innerxml"`
}

func extract(content []byte) (string, error) {
	h := html{}
	err := xml.NewDecoder(bytes.NewBuffer(content)).Decode(&h)
	if err != nil {
		fmt.Println("error", err)
		return "", err
	}

	return h.Body.Content, nil
}

func SendSignupEmail(name string, to string, from string, code string, subject string) error {

	t, exist := render.Render["SignupEmail"]
	html := ""

	if exist {
		var doc bytes.Buffer
		p := map[string]interface{}{"code": code, "subject": subject}
		t.Execute(&doc, p)
		html = doc.String()
	} else {
		return errors.New("no signup template")
	}

	password := config.DefaultConfig.APIKey
	value := "api:" + password
	auth := base64.StdEncoding.EncodeToString([]byte(value))
	text, err := extract([]byte(html))
	if err != nil {
		text = ""
	}
	ro := &grequests.RequestOptions{
		Data: map[string]string{
			"name":    name,
			"from":    from,
			"to":      to,
			"subject": subject,
			"html":    html,
			"text":    text,
		},
		Headers: map[string]string{"Authorization": "Basic " + auth},
	}

	mailServer := config.DefaultConfig.APIBaseUrl
	reqUrl := mailServer + "/mail/notice.hedwi.com"
	resp, err := grequests.Post(reqUrl, ro)
	// You can modify the request by passing an optional RequestOptions struct

	if err != nil {
		fmt.Println(err)
		return err
	}

	bs := resp.Bytes()
	res := map[string]interface{}{}
	err = json.Unmarshal(bs, &res)
	if err != nil {
		return err
	}
	if code, exist := res["code"].(float64); exist {
		if code == 0 {
			return nil
		} else {
			if message, exist := res["error"].(string); exist {
				return errors.New(message)
			} else {
				return errors.New("Errors when send email")
			}
		}
	}

	return nil
}

func SendForgotPasswordEmail(to string, from string, code string, subject string) error {

	t, exist := render.Render["ForgotPasswordEmail"]
	html := ""

	if exist {
		var doc bytes.Buffer
		p := map[string]interface{}{"subject": subject, "code": code}
		t.Execute(&doc, p)
		html = doc.String()
	} else {
		return errors.New("no forgotpassword template")
	}

	password := config.DefaultConfig.APIKey
	value := "api:" + password
	auth := base64.StdEncoding.EncodeToString([]byte(value))
	text, err := extract([]byte(html))
	if err != nil {
		text = ""
	}
	ro := &grequests.RequestOptions{
		Data: map[string]string{
			"from":    from,
			"to":      to,
			"subject": subject,
			"html":    html,
			"text":    text,
		},
		Headers: map[string]string{"Authorization": "Basic " + auth},
	}

	mailServer := config.DefaultConfig.APIBaseUrl
	reqUrl := mailServer + "/mail/notice.hedwi.com"
	resp, err := grequests.Post(reqUrl, ro)
	// You can modify the request by passing an optional RequestOptions struct

	if err != nil {
		fmt.Println(err)
		return err
	}

	bs := resp.Bytes()
	res := map[string]interface{}{}
	err = json.Unmarshal(bs, &res)
	if err != nil {
		return err
	}
	if code, exist := res["code"].(float64); exist {
		if code == 0 {
			return nil
		} else {
			if message, exist := res["error"].(string); exist {
				return errors.New(message)
			} else {
				return errors.New("Errors when send email")
			}
		}
	}

	return nil
}

func sendMail() {

	ssl := true
	host := config.DefaultConfig.SmtpServer
	username := "no-reply@notice.restdoc.com"
	//username = "hedwi@1024job.com"
	password := "1f6e26ad6f677d5415f5d8e9d15eb709"
	//password = "ujlqsgkgbfuyiaqd"
	auth := smtp.PlainAuth("", username, password, host)

	e := email.NewEmail()
	e.From = username
	//e.To = []string{"lxl360@qq.com"}
	e.To = []string{"839873980@qq.com"}
	e.Subject = "smtp test "
	e.Text = []byte("smtp test from hedwi")

	if ssl {
		err := e.SendWithTLS(host+":465", auth, &tls.Config{InsecureSkipVerify: false, ServerName: host})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("OK")
	} else {
		err := e.Send(host+":25", auth)
		if err != nil {
			fmt.Println(err)
		}
	}

}
