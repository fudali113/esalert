package alert

import (
	"bytes"
	"encoding/json"
	"net/http"
	"mylog"
	"util"
)

// HTTPAlert http 报警方式
type HTTPAlert struct {
	Url string
}

// Alert 发送http请求
func (httpAlert HTTPAlert) Alert(res map[string]interface{}) error {
	buffer := &bytes.Buffer{}
	bytes, _ := json.Marshal(res)
	buffer.Write(bytes)
	_, err := http.Post(httpAlert.Url, "application/josn", buffer)
	if err != nil {
		mylog.Error("http Alert 请求出错,", err)
	}
	return nil
}

// Create 。。。
func (HTTPAlert) GetAlerter(alertConfig map[string]interface{}) (alert interface{}, err error) {
	return HTTPAlert{Url: util.GetMapString(alertConfig, "url", "")}, nil
}

func (httpAlert HTTPAlert) GetTypes() []string {
	return []string{"http"}
}
