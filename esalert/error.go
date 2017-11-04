package esalert

type baseError struct {
	Message string
}

func (baseError baseError) Error() string {
	return baseError.Message
}

// ConfigError 配置错误
type ConfigError struct {
	Message string
}

func (configError ConfigError) Error() string {
	return configError.Message
}

// RequestError 请求错误
type RequestError struct {
	Message string
}

func (requestError RequestError) Error() string {
	return requestError.Message
}

type NotFoundError struct {
	Message string
}

func (notFoundError NotFoundError) Error() string {
	return notFoundError.Message
}

type RuleError struct {
	Message string
}

func (ruleError RuleError) Error() string {
	return ruleError.Message
}
