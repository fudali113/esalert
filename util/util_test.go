package util

import (
	"testing"
	bytes "bytes"
	"log"
	"html/template"
)

var (
	/*TestMail = Mail{
		Host:     "smtp.exmail.qq.com",
		Port:     "25",
		Username: "***",
		Password: "***",
		ReplyTo:  "fudali113@example.com",
	}
	SendTo = []string{"fuyi@23mofang.com"}*/
)

func TestMail_Send(t *testing.T) {
	/*Mail := TestMail
	err := Mail.Send(SendTo, "test", []byte("test llll 啪啪啪啪啪  \r\n  sfsdfdfdss saffasafsafafsbdhfhfd"))
	if err != nil {
		t.Error(err)
	}
	return*/
}

func TestTpl(t *testing.T)  {
	tpl := `{{ ._oo }} \n  {{ ._source }}`
	var temp *template.Template = template.New("ooo")
	temp, err := temp.Parse(tpl)
	if err != nil {
		t.Error(err)
	}
	buffer := bytes.NewBuffer([]byte{})
	temp.Execute(buffer, map[string]interface{}{"_oo": "oooooo", "_source": "pppppp"})
	log.Println(string(buffer.Bytes()))
}
