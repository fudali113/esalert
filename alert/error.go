package alert

type NotFoundError struct {
	Message string
}

func (notFoundError NotFoundError) Error() string {
	return notFoundError.Message
}

// ConfigError 配置错误
type ConfigError struct {
	Message string
}

func (configError ConfigError) Error() string {
	return configError.Message
}
