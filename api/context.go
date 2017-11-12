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
	return SplitAndSelect(c.r.RequestURI, "/", index)
}

func (c *Context) WriteJson(i interface{}) (err error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		c.w.Write([]byte(err.Error()))
		err = nil
		return
	}
	c.w.Header().Add("Content-Type", "application/json;charset=UTF-8")
	_, err = c.w.Write(bytes)
	return
}

func (c *Context) WriteString(s string) error {
	_, err := c.w.Write([]byte(s))
	return err
}

func (c *Context) WriteError(err error) error {
	return c.WriteString(err.Error())
}

func (c *Context) Write(bytes []byte) error {
	_, err := c.w.Write(bytes)
	return err
}

func (c *Context) Four04() {
	c.w.WriteHeader(404)
	c.WriteString("404 Not Found")
}
