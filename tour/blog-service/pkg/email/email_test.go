package email

import (
	"bytes"
	"strings"
	"testing"

	"gopkg.in/gomail.v2"
)

type mockDialer struct {
	msgs []*gomail.Message
	err  error
}

func (m *mockDialer) DialAndSend(msg ...*gomail.Message) error {
	m.msgs = append(m.msgs, msg...)
	return m.err
}

func TestSendMail_WithMockDialer(t *testing.T) {
	md := &mockDialer{}
	e := NewEmail(&SMTPInfo{
		Host:     "smtp.qq.com",
		Port:     465,
		UserName: "814594276@qq.com",
		Password: "gifrtnxipvidbdeb",
		From:     "814594276@qq.com",
		IsSSL:    true,
	}).WithDialer(md)

	err := e.SendMail([]string{"huangddab@163.com"}, "test subject", "<p>hello</p>")
	if err != nil {
		t.Fatalf("SendMail returned error: %v", err)
	}
	if len(md.msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(md.msgs))
	}
	msg := md.msgs[0]

	// 验证关键头部与主体
	if from := msg.GetHeader("From"); len(from) != 1 || from[0] != "814594276@qq.com" {
		t.Fatalf("unexpected From header: %v", from)
	}
	if to := msg.GetHeader("To"); len(to) != 1 || to[0] != "huangddab@163.com" {
		t.Fatalf("unexpected To header: %v", to)
	}
	if subj := msg.GetHeader("Subject"); len(subj) != 1 || subj[0] != "test subject" {
		t.Fatalf("unexpected Subject header: %v", subj)
	}
	var buf bytes.Buffer
	if _, err := msg.WriteTo(&buf); err != nil {
		t.Fatalf("failed to render message: %v", err)
	}
	if !strings.Contains(buf.String(), "<p>hello</p>") {
		t.Fatalf("body not found in message payload: %s", buf.String())
	}
}
