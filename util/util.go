package util

import (
	"bytes"
	"encoding/json"
	"net"
	"net/smtp"
	"log"
	"strings"
)

// ToBuffer 转换一个对象为byte[]
func ToBuffer(i interface{}) *bytes.Buffer {
	json, err := json.Marshal(i)
	if err != nil {
		log.Printf("转化json出现错误", err)
		return nil
	}
	return bytes.NewBuffer(json)
}

// QueryToJSON 转化query为json byte[]
func QueryToJSON(query interface{}) ([]byte, error) {
	query = CleanupMapValue(query)
	return json.Marshal(query)
}

// Mail 邮件实体
type Mail struct {
	Host, Port, Username, Password, From, ReplyTo string
}

// Send 发送邮件
// # BUG
// * 某些不确定情况下可能导致该函数阻塞，等待不确定的一段时间(可能很长，可能一会儿以后)报EOF错误,
//   跟断点猜测重要应该时golang源码没有超时机制和相关smtp服务器服务器有些问题(telnet客服端模拟也会出现发送HELO没有回复的情况)
func (mail Mail) Send(to []string, subject string, msg []byte) error {
	if mail.From == "" {
		mail.From = mail.Username
	}
	if mail.ReplyTo == "" {
		mail.ReplyTo = mail.From
	}
	// 如果msg中夹带空格，msg与header中间需要有一个分行smtp服务器才能识别内容
	msg = append([]byte("\r\n"), msg...)
	server := net.JoinHostPort(mail.Host, mail.Port)
	auth := smtp.PlainAuth("", mail.Username, mail.Password, mail.Host)
	from := []byte("From:" + mail.From)
	contentType := []byte("Content-Type: text/html")
	replyTo := []byte("Reply-To:" + mail.ReplyTo)
	sub := []byte("Subject:" + subject)
	return smtp.SendMail(server, auth, mail.From, to, bytes.Join([][]byte{contentType, from, replyTo, sub, msg}, []byte("\r\n")))
}

// MergeMap 融合两个map
func MergeMap(master, follow map[string]interface{}) map[string]interface{} {
	for k, v := range master {
		follow[k] = v
	}
	return follow
}

func BuildFileDir(dir string, name string) string {
	if strings.HasSuffix(dir, "/") {
		return dir + name
	}
	return dir + "/" + name
}
