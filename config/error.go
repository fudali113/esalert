package config

type ConfigError struct {
	Message string
}

func (error ConfigError) Error() string {
	return error.Message
}
