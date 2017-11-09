package alert

import (
	"bytes"
	"encoding/json"
	"net/http"
	"mylog"
	"config"
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
func (HTTPAlert) Create(config config.Config, alertConfig config.AlertConfig) (Alerter, error) {
	return HTTPAlert{Url: alertConfig.URL}, nil
}
