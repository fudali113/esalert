package esalert

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// EsRequest ealsticsearch请求对象
type EsRequest struct {
	host     string
	port     string
	name     string
	password string
	index    string
	body     map[string]interface{}
	request  *http.Request
}

// RunQuery 发起请求
func (er EsRequest) RunQuery() (res map[string]interface{}, err error) {
	if er.request == nil {
		var body io.Reader
		if er.body != nil {
			body = ToBuffer(er.body)
		}
		er.request, err = http.NewRequest("get", er.getURL(), body)
		if err != nil {
			return
		}
		er.request.SetBasicAuth(er.name, er.password)
		er.request.Header.Set("Content-Type", "Application/json")
	}
	if err != nil {
		return
	}
	response, err := http.DefaultClient.Do(er.request)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		err = RequestError{Message: string(body)}
		return
	}
	if err != nil {
		return
	}
	json.Unmarshal(body, &res)
	return
}

// getURL 获取查询地址
func (er EsRequest) getURL() string {
	return fmt.Sprintf("http://%s:%s/%s/_search", er.host, er.port, er.index)
}
