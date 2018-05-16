package util

import (
	"fmt"
	"github.com/gozh-io/gozh/module/configure"
	"gopkg.in/gomail.v2"
	"sync"
)

type EmailMessage struct {
	Headers     map[string][]string
	ContentType string
	Body        string
}

var (
	email      *gomail.Dialer
	email_once sync.Once
)

func Init() {
	email_once.Do(func() {
		conf := configure.GetConfigure()
		emailConf := conf.Email
		email = gomail.NewDialer(
			emailConf.Host,
			emailConf.Port,
			emailConf.UserName,
			emailConf.Password)
	})
}

func SendEmailTo(message EmailMessage) error {
	Init()
	m := gomail.NewMessage()
	m.SetHeaders(message.Headers)
	m.SetBody(message.ContentType, message.Body)
	if err := email.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendEmailCaptchaTo(recever, captcha string) error {
	Init()
	headers := make(map[string][]string)
	headers["From"] = []string{email.Username}
	headers["To"] = []string{recever}
	headers["Subject"] = []string{"gozh社区邮箱验证码"}
	body := fmt.Sprintf("您的邮箱验证码为: %v, 请勿告诉他人", captcha)
	fmt.Println(body, recever)
	em := EmailMessage{Headers: headers, Body: body, ContentType: "text/html"}
	if err := SendEmailTo(em); err != nil {
		fmt.Println("发送失败")
		return fmt.Errorf("发送验证码到邮箱%v失败", recever)
	}
	fmt.Println("发送成功")
	return nil
}
