package alert

import (
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
func (MailAlert) GetAlerter(alertConfig map[string]interface{}) (alert interface{}, err error) {
	alert = MailAlert{}
	var template *template.Template
	content := util.GetMapString(alertConfig, "content", "")
	tplFile := util.GetMapString(alertConfig, "tpl_file", "")
	if tplFile == "" {
		if util.GetMapString(alertConfig, "content", "") == "" {
			err = ConfigError{Message: "tpl_file || content must exists by Mail"}
			return
		}
		template, err = template.Parse(content)
		if err != nil {
			return
		}
	} else {
		template, err = template.ParseFiles(tplFile)
		if err != nil {
			return
		}
		bytes, err := ioutil.ReadFile(tplFile)
		if err == nil {
			content = string(bytes)
		}
	}
	alert = MailAlert{
		Mail: util.Mail{
			Host:     util.GetMapString(alertConfig, "smtp_host", ""),
			Port:     util.GetMapString(alertConfig, "smtp_port", ""),
			Username: util.GetMapString(alertConfig, "username", ""),
			Password: util.GetMapString(alertConfig, "password", ""),
			From:     util.GetMapString(alertConfig, "from_addr", ""),
			ReplyTo:  util.GetMapString(alertConfig, "reply_to", ""),
		},
		To:       util.GetMapStringSlice(alertConfig, "send_to", []string{}),
		Subject:  util.GetMapString(alertConfig, "subject", ""),
		Content:  content,
		template: template,
	}
	return
}

func (mailAlert MailAlert) GetTypes() []string {
	return []string{"mail"}
}
