package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"mylog"
)

type ConfigDirInfo struct {
	Dir        string
	RuleName   string
	ConfigName string
}

// Config 全部配置
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Rules    []RuleConfig
	Mail     MailConfig
	ApiInfo  ApiConfig `yaml:"api"`
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
	Index    string
	Body     map[string]interface{}
	Script   string
	Test     TestRuleData
	Interval Time
	Alerts   []AlertConfig
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

// AlertConfig 报警配置
type AlertConfig struct {
	Type string
	URL  string `yaml:"Url"`
	Mail MailConfig
}

// MailConfig 邮箱信息配置
type MailConfig struct {
	Username string
	Password string
	SMTPHost string   `yaml:"smtp_host"`
	SMTPPort string   `yaml:"smtp_port"`
	SendTo   []string `yaml:"send_to"`
	FromAddr string   `yaml:"from_addr"`
	ReplyTo  string   `yaml:"reply_to"`
	TPLFile  string   `yaml:"tpl_file"`
	Content  string
	Subject  string
}

var (
	OriginConfig *Config
)

// IntiConfig 根据配置文件路径加载配置
func IntiConfig(configDirInfo ConfigDirInfo) (config *Config, err error) {
	bytes, err := ioutil.ReadFile(buildFileDir(configDirInfo.Dir, configDirInfo.ConfigName))
	if err != nil {
		return nil, err
	}
	config = &Config{
		Host:     "localhsot",
		Port:     "9200",
		Username: "elastic",
		Password: "changeme",
		Rules:    []RuleConfig{},
	}
	// 保留原始配置
	defer func() {
		OriginConfig = config
	}()
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}
	ruleDir := buildFileDir(configDirInfo.Dir, configDirInfo.RuleName)
	rulesFile, err := ioutil.ReadDir(ruleDir)
	if err != nil {
		mylog.Info("打开规则文件"+ruleDir+"夹出错", err)
		return config, nil
	}
	rules := []RuleConfig{}
	for _, ruleFile := range rulesFile {
		if !ruleFile.IsDir() {
			rule, err := GetRuleConfig(buildFileDir(ruleDir, ruleFile.Name()))
			if err != nil {
				mylog.Error(err)
				continue
			}
			rules = append(rules, rule)
		}
	}
	config.Rules = append(config.Rules, rules...)
	if len(rules) == 0 {
		return config, ConfigError{"规则数不能为0"}
	}
	return config, nil
}

func GetRuleConfig(ruleDir string) (RuleConfig, error) {
	rule := RuleConfig{}
	bytes, err := ioutil.ReadFile(ruleDir)
	if err != nil {
		return RuleConfig{}, err
	}
	yaml.Unmarshal(bytes, &rule)
	return rule, nil
}

func buildFileDir(dir string, name string) string {
	if strings.HasSuffix(dir, "/") {
		return dir + name
	}
	return dir + "/" + name
}
