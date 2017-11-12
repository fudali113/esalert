package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"mylog"
	"util"
	"fmt"
)

type ConfigDirInfo struct {
	Dir        string
	RuleName   string
	ConfigName string
}

// Config 全部配置
type Config struct {
	Storage map[string]interface{}
	Rules   []RuleConfig
	Alert   map[string]interface{}
	ApiInfo ApiConfig `yaml:"api"`
	Test    bool
}

// Valid 验证config信息
// TODO 更详细的验证信息
func (config Config) Valid() error {
	ruleNameMap := map[string]int{}
	for _, rule := range config.Rules {
		if _, ok := ruleNameMap[rule.Name]; ok {
			return fmt.Errorf("rule name 必须唯一，出现了两个相同name 的 rule: %s", rule.Name)
		}
		ruleNameMap[rule.Name] = 1
	}
	return nil
}

type ApiConfig struct {
	Enable    bool
	Port      string
	BasicAuth BasicAuthConfig `yaml:"basic_auth"`
}

type BasicAuthConfig struct {
	Enable   bool
	Username string
	Password string
}

// RuleConfig 规则配置
type RuleConfig struct {
	Name     string
	Status   int
	Storage  map[string]interface{}
	Script   string
	Test     TestRuleData
	Interval Time
	Alerts   []map[string]interface{}
}

type TestRuleData struct {
	Data   string
	Should bool
}

// Time 语义相关的时间
// such as ["day:1","hour":2,"second": 30] = 1天2小时30秒
type Time map[string]int32

// GetSecond 根据时间信息获取秒数
func (t Time) GetSecond() int32 {
	var second int32
	for k, v := range t {
		switch k {
		case "y", "Y":
			second += 60 * 60 * 24 * 365 * v
		case "M":
			second += 60 * 60 * 24 * 30 * v
		case "d", "D":
			second += 60 * 60 * 24 * v
		case "h", "H":
			second += 60 * 60 * v
		case "m":
			second += 60 * v
		case "s", "S":
			second += v
		}
	}
	return second
}

var (
	OriginConfig *Config
)

// IntiConfig 根据配置文件路径加载配置
func IntiConfig(configDirInfo ConfigDirInfo) (config *Config, err error) {
	bytes, err := ioutil.ReadFile(util.BuildFileDir(configDirInfo.Dir, configDirInfo.ConfigName))
	if err != nil {
		return nil, err
	}
	config = &Config{
		Storage: map[string]interface{}{
			"host":     "localhsot",
			"port":     "9200",
			"username": "elastic",
			"password": "changeme",
		},
		Rules: []RuleConfig{},
	}
	util.CleanupMapValue(config.Storage)
	util.CleanupMapValue(config.Alert)
	// 保留原始配置
	defer func() {
		OriginConfig = config
	}()
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}
	ruleDir := util.BuildFileDir(configDirInfo.Dir, configDirInfo.RuleName)
	rulesFile, err := ioutil.ReadDir(ruleDir)
	if err != nil {
		mylog.Info("打开规则文件"+ruleDir+"夹出错", err)
		return config, nil
	}
	rules := []RuleConfig{}
	for _, ruleFile := range rulesFile {
		if !ruleFile.IsDir() {
			rule, err := GetRuleConfig(util.BuildFileDir(ruleDir, ruleFile.Name()))
			if err != nil {
				mylog.Error(err)
				continue
			}
			rules = append(rules, rule)
		}
	}
	config.Rules = append(config.Rules, rules...)
	for _, rule := range config.Rules {
		rule.Storage = util.CleanupStringMap(rule.Storage)
		for i, alert := range rule.Alerts {
			rule.Alerts[i] = util.CleanupStringMap(alert)
		}
	}
	if len(rules) == 0 {
		return config, ConfigError{"规则数不能为0"}
	}
	return config, config.Valid()
}

func GetRuleConfig(ruleDir string) (RuleConfig, error) {
	rule := RuleConfig{Storage: map[string]interface{}{}, Alerts: []map[string]interface{}{}}
	bytes, err := ioutil.ReadFile(ruleDir)
	if err != nil {
		return RuleConfig{}, err
	}
	yaml.Unmarshal(bytes, &rule)
	return rule, nil
}
