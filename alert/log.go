package alert

import (
	"mylog"
)

// LogAlert 打印日志报警方式，默认报警方式，当没有任何报警方式时，自动添加该报警方式
type LogAlert struct {
}

// Alert 打印日志记录
func (LogAlert) Alert(res map[string]interface{}) error {
	mylog.Info(res)
	return nil
}

// Create 。。。
func (LogAlert) GetAlerter(map[string]interface{}) (interface{}, error) {
	return LogAlert{}, nil
}

func (LogAlert) GetTypes() []string {
	return []string{"log", "default"}
}
