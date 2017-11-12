package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"util"
)

// EsStorage ealsticsearch请求对象
type EsStorage struct {
	Host     string
	Port     string
	Username string
	Password string
	Index    string
	Body     interface{}
	request  *http.Request
}

// GetData 发起请求
func (es EsStorage) GetData() (res map[string]interface{}, err error) {
	if es.request == nil {
		var body io.Reader
		if es.Body != nil {
			body = util.ToBuffer(es.Body)
		}
		es.request, err = http.NewRequest("get", es.getURL(), body)
		if err != nil {
			return
		}
		es.request.SetBasicAuth(es.Username, es.Password)
		es.request.Header.Set("Content-Type", "Application/json")
	}
	if err != nil {
		return
	}
	response, err := http.DefaultClient.Do(es.request)
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

// 改储存的名字
func (EsStorage) GetTypes() []string {
	return []string{"es", "elasticsearch", "default"}
}

func (EsStorage) GetStorage(config map[string]interface{}) (es interface{}, err error) {
	es = EsStorage{}
	defer func() {
		if e := recover(); e != nil {
			if e, ok := e.(error); ok {
				err = e
				return
			}
		}
	}()
	err = valid(config)
	if err != nil {
		err = fmt.Errorf("EsStorage 配置文件出错， error: %s", err.Error())
		return
	}
	es = EsStorage{
		Host:     util.GetMapString(config, "host", "localshot"),
		Port:     util.GetMapString(config, "port", "9200"),
		Username: util.GetMapString(config, "username", "elastic"),
		Password: util.GetMapString(config, "password", "localshot"),
		Index:    util.GetMapString(config, "index", "localshot"),
		Body:     util.CleanupMapValue(config["body"]),
	}
	return
}

// valid 验证参数
func valid(config map[string]interface{}) error {
	err := util.AssertMapHas(config, "host")
	err += util.AssertMapHas(config, "port")
	err += util.AssertMapHas(config, "username")
	err += util.AssertMapHas(config, "password")
	err += util.AssertMapHas(config, "index")
	if err != "" {
		return fmt.Errorf(err)
	}
	return nil
}

// getURL 获取查询地址
func (er EsStorage) getURL() string {
	return fmt.Sprintf("http://%s:%s/%s/_search", er.Host, er.Port, er.Index)
}
