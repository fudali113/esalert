package api

import (
	"net/http"
	"encoding/json"
)

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func (c *Context) URISubSelect(index int) string {
	return SplitAndSelect(c.r.RequestURI, "/", -1)
}

func (c *Context) WriteJson(i interface{}) error {
	bytes, err := json.Marshal(i)
	if err != nil {
		return err
	}
	c.w.Header().Add("Content-Type", "application/json;charset=UTF-8")
	_, err = c.w.Write(bytes)
	return err
}

func (c *Context) WriteString(s string) error {
	_, err := c.w.Write([]byte(s))
	return err
}

func (c *Context) Write(bytes []byte) error {
	_, err := c.w.Write(bytes)
	return err
}
