package alert

import (
	"config"
	"bytes"
	"io/ioutil"
	"util"
	"html/template"
)

// MailAlert 发送邮件报警方式
type MailAlert struct {
	Mail     util.Mail
	To       []string
	Subject  string
	Content  string
	template *template.Template
}

// Alert 发送邮件
func (mailAlert MailAlert) Alert(res map[string]interface{}) error {
	buffer := bytes.NewBuffer([]byte{})
	mailAlert.template.Execute(buffer, res)
	return mailAlert.Mail.Send(mailAlert.To, mailAlert.Subject, buffer.Bytes())
}

// Create 。。。
func (MailAlert) Create(config config.Config, alertConfig config.AlertConfig) (Alerter, error) {
	mail := merge(alertConfig.Mail, config.Mail)
	var err error
	var template *template.Template
	content := ""
	if mail.TPLFile == "" {
		if mail.Content == "" {
			return nil, ConfigError{Message: "tpl_file || content must exists by Mail"}
		}
		content = mail.Content
		template, err = template.Parse(mail.Content)
		if err != nil {
			return nil, err
		}
	} else {
		template, err = template.ParseFiles(mail.TPLFile)
		if err != nil {
			return nil, err
		}
		bytes, err := ioutil.ReadFile(mail.TPLFile)
		if err == nil {
			content = string(bytes)
		}
	}
	return MailAlert{
		Mail: util.Mail{
			Host:     mail.SMTPHost,
			Port:     mail.SMTPPort,
			Username: mail.Username,
			Password: mail.Password,
			From:     mail.FromAddr,
			ReplyTo:  mail.ReplyTo,
		},
		To:       mail.SendTo,
		Subject:  mail.Subject,
		Content:  content,
		template: template,
	}, nil
}

func merge(base config.MailConfig, config config.MailConfig) config.MailConfig {
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
