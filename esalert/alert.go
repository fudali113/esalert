package esalert

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var registerMap = map[string]AlerterCreater{}

// Register 注册一个AlerterCreater
func Register(t string, creater AlerterCreater) {
	t = strings.ToLower(t)
	_, ok := registerMap[t]
	if ok {
		log.Println("WARN ", "已经存在一个该type的creater")
	}
	registerMap[t] = creater
}

// GetCreater 获取一个AlerterCreater
func GetCreater(t string) (AlerterCreater, error) {
	t = strings.ToLower(t)
	v, ok := registerMap[t]
	if !ok {
		return nil, NotFoundError{}
	}
	return v, nil
}

// CreateAlerter 根据type生成一个Alerter
func CreateAlerter(t string, config Config, alertConfig AlertConfig) (Alerter, error) {
	creater, err := GetCreater(t)
	if err != nil {
		return nil, err
	}
	return creater.Create(config, alertConfig)
}

func init() {
	Register("log", LogAlert{})
	Register("http", HTTPAlert{})
	Register("mail", MailAlert{})
}

// Alerter 报警方式处理接口
type Alerter interface {
	// Alert 根据结果报警
	Alert(res map[string]interface{}) error
}

// AlerterCreater 产生者
// 将产生逻辑规定到Alert的出生地
type AlerterCreater interface {
	// Create 根据config生成一个alert
	Create(config Config, alertConfig AlertConfig) (Alerter, error)
}

// LogAlert 打印日志报警方式，默认报警方式，当没有任何报警方式时，自动添加该报警方式
type LogAlert struct {
}

// Alert 打印日志记录
func (LogAlert) Alert(res map[string]interface{}) error {
	log.Println(res)
	return nil
}

// Create 。。。
func (LogAlert) Create(config Config, alertConfig AlertConfig) (Alerter, error) {
	return LogAlert{}, nil
}

// HTTPAlert http 报警方式
type HTTPAlert struct {
	url string
}

// Alert 发送http请求
func (httpAlert HTTPAlert) Alert(res map[string]interface{}) error {
	buffer := &bytes.Buffer{}
	bytes, _ := json.Marshal(res)
	buffer.Write(bytes)
	_, err := http.Post(httpAlert.url, "application/josn", buffer)
	if err != nil {
		log.Print("http Alert 请求出错,", err)
	}
	return nil
}

// Create 。。。
func (HTTPAlert) Create(config Config, alertConfig AlertConfig) (Alerter, error) {
	return HTTPAlert{url: alertConfig.URL}, nil
}

// MailAlert 发送邮件报警方式
type MailAlert struct {
	mail     Mail
	to       []string
	subject  string
	template *template.Template
}

// Alert 发送邮件
func (mailAlert MailAlert) Alert(res map[string]interface{}) error {
	buffer := bytes.NewBuffer([]byte{})
	mailAlert.template.Execute(buffer, res)
	return mailAlert.mail.Send(mailAlert.to, mailAlert.subject, buffer.Bytes())
}

// Create 。。。
func (MailAlert) Create(config Config, alertConfig AlertConfig) (Alerter, error) {
	mail := merge(alertConfig.Mail, config.Mail)
	var err error
	var template *template.Template
	if mail.TPLFile == "" {
		if mail.Content == "" {
			return nil, ConfigError{Message: "tpl_file || content must exists by mail"}
		}
		template, err = template.Parse(mail.Content)
		if err != nil {
			return nil, err
		}
	} else {
		template, err = template.ParseFiles(mail.TPLFile)
		if err != nil {
			return nil, err
		}
	}
	return MailAlert{
		mail: Mail{
			Host:     mail.SMTPHost,
			Port:     mail.SMTPPort,
			Username: mail.Username,
			Password: mail.Password,
			From:     mail.FromAddr,
			ReplyTo:  mail.ReplyTo,
		},
		to:       mail.SendTo,
		subject:  mail.Subject,
		template: template,
	}, nil
}

func merge(base MailConfig, config MailConfig) MailConfig {
	if base.Username == "" {
		base.Username = config.Username
	}
	if base.Password == "" {
		base.Password = config.Password
	}
	if base.SMTPHost == "" {
		base.SMTPHost = config.SMTPHost
	}
	if base.SMTPPort == "" {
		base.SMTPPort = config.SMTPPort
	}
	if len(base.SendTo) == 0 {
		base.SendTo = config.SendTo
	}
	if base.FromAddr == "" {
		base.FromAddr = config.FromAddr
	}
	return base
}
