package esalert

import (
	"io/ioutil"
	"time"

	"context"

	"gopkg.in/yaml.v2"
)

// RuleContext 保存各个rule允许的状态
var RuleContext = map[string]*RuleRunContext{}

// RuleRunContext rule运行的上下文
// TODO 保存更多rule运行相关状态
type RuleRunContext struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// Stop 停止某个rule
func (rrc *RuleRunContext) Stop() {
	rrc.cancelFunc()
}

// IntiConfig 根据配置文件路径加载配置
func IntiConfig(configDir string) (*Config, error) {
	bytes, err := ioutil.ReadFile(configDir)
	if err != nil {
		return nil, err
	}
	config := &Config{
		Host:     "localhsot",
		Port:     "9200",
		Username: "elastic",
		Password: "changeme",
		Rules:    []RuleConfig{},
	}
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Run 启动配置参数
func Run(config Config) error {
	if len(config.Rules) == 0 {
		return ConfigError{Message: "rules不能为空"}
	}
	rules := []rule{}
	for _, rule := range config.Rules {
		rules = append(rules, sampleRule{
			name:      rule.Name,
			esRequest: getEsRequest(config, rule),
			tick:      time.NewTicker(time.Duration(rule.Interval.GetSecond()) * time.Second),
			script:    rule.Script,
			alerter:   getAlerts(config, rule.Alerts),
		})
	}
	for _, rule := range rules {
		ctx, cancelFunc := context.WithCancel(context.Background())
		rule.run(ctx)
		RuleContext[rule.Name()] = &RuleRunContext{ctx: ctx, cancelFunc: cancelFunc}
	}
	return nil
}

func getAlerts(config Config, alertConfigs []AlertConfig) []Alerter {
	alerterList := make([]Alerter, 0, len(alertConfigs))
	for _, alertConfig := range alertConfigs {
		alerter, err := CreateAlerter(alertConfig.Type, config, alertConfig)
		if err != nil {
			panic(err)
		}
		alerterList = append(alerterList, alerter)
	}
	if len(alerterList) == 0 {
		alerterList = append(alerterList, LogAlert{})
	}
	return alerterList
}

func getEsRequest(config Config, rule RuleConfig) EsRequest {
	return EsRequest{
		host:     config.Host,
		port:     config.Port,
		name:     config.Username,
		password: config.Password,
		index:    rule.Index,
		body:     cleanupStringMap(rule.Body),
	}
}
