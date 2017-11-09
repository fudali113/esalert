package source

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"util"
)

// EsRequest ealsticsearch请求对象
type EsRequest struct {
	Host     string
	Port     string
	Name     string
	Password string
	Index    string
	Body     map[string]interface{}
	request  *http.Request
}

// RunQuery 发起请求
func (er EsRequest) RunQuery() (res map[string]interface{}, err error) {
	if er.request == nil {
		var body io.Reader
		if er.Body != nil {
			body = util.ToBuffer(er.Body)
		}
		er.request, err = http.NewRequest("get", er.getURL(), body)
		if err != nil {
			return
		}
		er.request.SetBasicAuth(er.Name, er.Password)
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
	return fmt.Sprintf("http://%s:%s/%s/_search", er.Host, er.Port, er.Index)
}
