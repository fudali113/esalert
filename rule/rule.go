package rule

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mylog"
	"javascript"
	"alert"
	"storage"
	"config"
)

type Rule interface {
	GetName() string
	Run(ctx context.Context)
}

func CreateRule(con config.RuleConfig) Rule {
	return nil
}

type SampleRule struct {
	Name    string
	Storage storage.Storage
	Tick    *time.Ticker
	Time    int32
	Script  string
	Alerter []alert.Alerter
}

func (rule SampleRule) GetName() string {
	return rule.Name
}

func (rule SampleRule) Run(ctx context.Context) {
	runFunc := func() {
		mylog.Info(fmt.Sprintf("rule: %s runing", rule.Name))
		res, err := rule.Storage.GetData()
		if err != nil {
			mylog.Error(err)
			return
		}
		resByte, _ := json.Marshal(res)
		mylog.Info(string(resByte))
		alert, err := NeedAlert(res, rule.Script)
		if err != nil {
			mylog.Error("运行判断是否需要报警脚本出错，", err)
		}
		if alert {
			mylog.Info("经过脚本运算需要发送提醒")
			res["_rule_info"] = rule
			for _, alerter := range rule.Alerter {
				err := alerter.Alert(res)
				if err != nil {
					mylog.Error(err)
					// continue
				}
			}
		}
		mylog.Info(fmt.Sprintf("rule: %s Run success", rule.Name))
	}
	go func() {
		runFunc()
	}()
	go func() {
		for {
			select {
			case <-rule.Tick.C:
				runFunc()
			case <-ctx.Done():
				mylog.Info(fmt.Sprintf("rule: %s stoped", rule.Name))
				return
			}
		}
	}()
}

// NeedAlert 根据结果判断是否需要发送报警信息
func NeedAlert(res map[string]interface{}, script string) (alert bool, err error) {
	vm := javascript.GetVM()
	vm.Set("result", res)
	value, err := vm.Run(script)
	if err != nil {
		return
	}
	return value.ToBoolean()
}
