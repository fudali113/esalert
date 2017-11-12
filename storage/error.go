package storage

// RequestError 请求错误
type RequestError struct {
	Message string
}

func (requestError RequestError) Error() string {
	return requestError.Message
}
