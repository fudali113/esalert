package esalert

// Config 全部配置
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Rules    []RuleConfig
	Mail     MailConfig
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

// RuleConfig 规则配置
type RuleConfig struct {
	Name     string
	Index    string
	Body     map[string]interface{}
	Script   string
	Interval Time
	Alerts   []AlertConfig
}

// Time 语义相关的时间
// such as ["day:1","hour":2,"second": 30] = 1天2小时30秒
type Time map[string]int32

// GetSecond 根据时间信息获取秒数
func (t Time) GetSecond() int32 {
	var second int32
	for k, v := range t {
		switch k {
		case "year":
			second += 60 * 60 * 24 * 365 * v
		case "month":
			second += 60 * 60 * 24 * 30 * v
		case "day":
			second += 60 * 60 * 24 * v
		case "hour":
			second += 60 * 60 * v
		case "minute":
			second += 60 * v
		case "second":
			second += v
		}
	}
	return second
}

// AlertConfig 报警配置
type AlertConfig struct {
	Type string
	URL  string `yaml:"url"`
	Mail MailConfig
}
