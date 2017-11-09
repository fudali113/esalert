package alert

import (
	"strings"
	"mylog"
	"config"
)

var registerMap = map[string]AlerterCreater{}

// Register 注册一个AlerterCreater
func Register(t string, creater AlerterCreater) {
	t = strings.ToLower(t)
	_, ok := registerMap[t]
	if ok {
		mylog.Warn("已经存在一个该type的creater")
	}
	registerMap[t] = creater
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
func CreateAlerter(t string, config config.Config, alertConfig config.AlertConfig) (Alerter, error) {
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
	Create(config config.Config, alertConfig config.AlertConfig) (Alerter, error)
}
