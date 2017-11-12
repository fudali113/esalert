package rule

import (
	"config"
	"time"
	"alert"
	"storage"
	"context"
	"fmt"
	"util"
)

// RuleContextMap 保存各个rule允许的状态
var RuleContextMap = map[string]*RuleContext{}

// RuleContext rule运行的上下文
type RuleContext struct {
	Status     int
	ctx        context.Context
	cancelFunc context.CancelFunc
	rule       Rule
	RuleConfig config.RuleConfig
}

const (
	STOP  = iota
	START
)

// Stop 停止某个rule
func (rc *RuleContext) Stop() error {
	rc.cancelFunc()
	rc.Status = STOP
	return nil
}

// Start 启动某个rule
func (rc *RuleContext) Start() error {
	if rc.RuleConfig.Name == "" {
		return fmt.Errorf("原始配置文件不能为空")
	}
	rule, err := GetRule(rc.RuleConfig)
	if err != nil {
		return err
	}
	rc.rule = rule
	rc.ctx, rc.cancelFunc = context.WithCancel(context.Background())
	rc.rule.Run(rc.ctx)
	rc.Status = START
	return nil
}

// Restart 如果已经启动，重新启动；如果没有启动，启动该rule
func (rc *RuleContext) Restart() error {
	if rc.Status == START {
		err := rc.Stop()
		if err != nil {
			return err
		}
	}
	return rc.Start()
}

// Run 启动配置参数
func Run(con config.Config) error {
	errs := []error{}
	for _, rule := range con.Rules {
		_, err := RunRule(rule)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	if con.Test && len(errs) > 0 {
		// TODO 将错误信息聚合起来
		return fmt.Errorf("配置文件有错！")
	}
	return nil
}

// RunRule 根据rule配置运行rule并保存在context中，会自动合并总配置文件的配置
func RunRule(ruleConfig config.RuleConfig) (rc *RuleContext, err error) {
	if _, ok := RuleContextMap[ruleConfig.Name]; ok {
		err = fmt.Errorf("RuleConfig name " + ruleConfig.Name + "已经存在")
		return
	}
	ruleConfig = MergeConfig(ruleConfig)
	if err != nil {
		return
	}
	rc = &RuleContext{RuleConfig: ruleConfig}
	if ruleConfig.Status == START {
		rc.Start()
	}
	RuleContextMap[ruleConfig.Name] = rc
	return
}

// GetRule 根据配置获取rule运行时
func GetRule(ruleConfig config.RuleConfig) (rule Rule, err error) {
	s, err := storage.GetStorage(ruleConfig.Storage)
	if err != nil {
		err = fmt.Errorf("rule %s : %s", ruleConfig.Name, err.Error())
		return
	}
	rule = SampleRule{
		Name:    ruleConfig.Name,
		Storage: s,
		Tick:    time.NewTicker(time.Duration(ruleConfig.Interval.GetSecond()) * time.Second),
		Time:    ruleConfig.Interval.GetSecond(),
		Script:  ruleConfig.Script,
		Alerter: getAlerts(ruleConfig.Alerts),
	}
	return
}

// MergeConfig 合并原始信息
func MergeConfig(rule config.RuleConfig) config.RuleConfig {
	rule.Storage = util.MergeNewMap(rule.Storage, config.OriginConfig.Storage)
	for i, alert := range rule.Alerts {
		rule.Alerts[i] = util.MergeNewMap(alert, config.OriginConfig.Alert)
	}
	return rule
}

// getAlerts 获取报警
func getAlerts(alertConfigs []map[string]interface{}) []alert.Alerter {
	alerterList := make([]alert.Alerter, 0, len(alertConfigs))
	for _, alertConfig := range alertConfigs {
		alerter, err := alert.CreateAlerter(alertConfig)
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
