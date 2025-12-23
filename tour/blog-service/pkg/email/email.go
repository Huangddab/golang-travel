package email

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type Email struct {
	*SMTPInfo
	dialer dialer
}

type SMTPInfo struct {
	Host     string
	Port     int
	IsSSL    bool
	UserName string
	Password string
	From     string
}

// 抽象拨号器便于测试注入
type dialer interface {
	DialAndSend(msg ...*gomail.Message) error
}

// 创建消息实例
func NewEmail(info *SMTPInfo) *Email {
	return &Email{SMTPInfo: info}
}

// WithDialer 允许自定义拨号器（用于单元测试或自定义发送策略）
func (e *Email) WithDialer(d dialer) *Email {
	e.dialer = d
	return e
}

func (e *Email) SendMail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.From)     // 发件人
	m.SetHeader("To", to...)        // 收件人
	m.SetHeader("Subject", subject) // 邮件主题
	m.SetBody("text/html", body)    // 邮件正文

	// 选择拨号器（优先使用注入的，便于测试）
	d := e.dialer
	if d == nil {
		dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.Password)
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}
		d = dialer
	}

	return d.DialAndSend(m)
}
