package alert

import (
	"strings"
	"mylog"
	"fmt"
)

var registerMap = map[string]AlerterCreater{}

// Register 注册一个AlerterCreater
func Register(creater AlerterCreater) {
	names := creater.GetTypes()
	alerter, _ := creater.GetAlerter(map[string]interface{}{})
	if alerter == nil {
		panic(fmt.Errorf("config 参数为空时必须返回一个不为nil的Alerter对象"))
	}
	if _, ok := alerter.(Alerter); !ok {
		panic(fmt.Errorf("返回的类型必须是Alerter类型"))
	}
	for _, name := range names {
		name = strings.ToLower(name)
		if _, ok := registerMap[name]; ok {
			mylog.Warn("已经存在一个该type的alerter creater, type: " + name)
		}
		registerMap[name] = creater
	}
}

// GetCreater 获取一个AlerterCreater
func GetCreater(t string) (AlerterCreater, error) {
	t = strings.ToLower(t)
	v, ok := registerMap[t]
	if !ok {
		return nil, NotFoundError{"not found"}
	}
	return v, nil
}

// CreateAlerter 根据type生成一个Alerter
func CreateAlerter(alertConfig map[string]interface{}) (Alerter, error) {
	alertType, ok := alertConfig["_type"]
	if !ok {
		alertType = "default"
	}
	creater, err := GetCreater(alertType.(string))
	if err != nil {
		return nil, err
	}
	aleter, err := creater.GetAlerter(alertConfig)
	if err != nil {
		return nil, err
	}
	return aleter.(Alerter), nil
}

func init() {
	Register(LogAlert{})
	Register(HTTPAlert{})
	Register(MailAlert{})
}

// Alerter 报警方式处理接口
type Alerter interface {
	// Alert 根据结果报警
	Alert(res map[string]interface{}) error
}

// AlerterCreater 产生者
// 将产生逻辑规定到Alert的出生地
type AlerterCreater interface {
	GetTypes() []string
	// Create 根据config生成一个alert
	GetAlerter(alertConfig map[string]interface{}) (interface{}, error)
}
