package esalert

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/robertkrimen/otto"
)

type rule interface {
	Name() string
	run(ctx context.Context)
}

type sampleRule struct {
	name      string
	esRequest EsRequest
	tick      *time.Ticker
	time      int32
	script    string
	alerter   []Alerter
}

func (rule sampleRule) Name() string {
	return rule.name
}

func (rule sampleRule) run(ctx context.Context) {
	runFunc := func() {
		log.Println("INFO ", fmt.Sprintf("rule: %s runing", rule.name))
		res, err := rule.esRequest.RunQuery()
		if err != nil {
			log.Println("ERROR ", err)
			return
		}
		alert, err := needAlert(res, rule.script)
		if err != nil {
			log.Println("INFO  运行判断是否需要报警脚本出错，", err)
		}
		if alert {
			for _, alerter := range rule.alerter {
				err := alerter.Alert(res)
				if err != nil {
					log.Println("ERROR", err)
					// continue
				}
			}
		}
		log.Println("INFO ", fmt.Sprintf("rule: %s run success", rule.name))
	}
	go func() {
		runFunc()
	}()
	go func() {
		for {
			select {
			case <-rule.tick.C:
				runFunc()
			case <-ctx.Done():
				log.Println("INFO ", fmt.Sprintf("rule: %s stoped", rule.name))
				break
			}
		}
	}()
}

// needAlert 根据结果判断是否需要发送报警信息
func needAlert(res map[string]interface{}, script string) (alert bool, err error) {
	vm := otto.New()
	vm.Set("result", res)
	value, err := vm.Run(script)
	if err != nil {
		return
	}
	return value.ToBoolean()
}
