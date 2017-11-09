package rule

import (
	"config"
	"time"
	"alert"
	"source"
	"util"
	"context"
)

// RuleContextMap 保存各个rule允许的状态
var RuleContextMap = map[string]*RuleContext{}

// RuleContext rule运行的上下文
type RuleContext struct {
	Status     int  `json:"status"`
	ctx        context.Context
	cancelFunc context.CancelFunc
	Rule       Rule `json:"rule"`
}

// Stop 停止某个rule
func (rc *RuleContext) Stop() {
	rc.cancelFunc()
	rc.Status = 0
}

func (rc *RuleContext) Start() {
	rc.ctx, rc.cancelFunc = context.WithCancel(context.Background())
	rc.Rule.Run(rc.ctx)
	rc.Status = 1
}

// Run 启动配置参数
func Run(con config.Config) error {
	for _, rule := range con.Rules {
		runRule := SampleRule{
			Name:      rule.Name,
			EsRequest: getEsRequest(*config.OriginConfig, rule),
			Tick:      time.NewTicker(time.Duration(rule.Interval.GetSecond()) * time.Second),
			Time:      rule.Interval.GetSecond(),
			Script:    rule.Script,
			Alerter:   getAlerts(*config.OriginConfig, rule.Alerts),
		}
		RuleContextMap[rule.Name] = &RuleContext{Rule: runRule}
	}
	for _, rule := range RuleContextMap {
		rule.Start()
	}
	return nil
}

func getAlerts(config config.Config, alertConfigs []config.AlertConfig) []alert.Alerter {
	alerterList := make([]alert.Alerter, 0, len(alertConfigs))
	for _, alertConfig := range alertConfigs {
		alerter, err := alert.CreateAlerter(alertConfig.Type, config, alertConfig)
		if err != nil {
			panic(err)
		}
		alerterList = append(alerterList, alerter)
	}
	if len(alerterList) == 0 {
		alerterList = append(alerterList, alert.LogAlert{})
	}
	return alerterList
}

func getEsRequest(config config.Config, rule config.RuleConfig) source.EsRequest {
	return source.EsRequest{
		Host:     config.Host,
		Port:     config.Port,
		Name:     config.Username,
		Password: config.Password,
		Index:    rule.Index,
		Body:     util.CleanupStringMap(rule.Body),
	}
}
