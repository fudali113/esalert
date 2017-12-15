package api

import (
	"net/http"
	"strings"
	"regexp"
)

var defaultMedias = []media{}

type restFuncMap map[string]http.Handler

func (rest restFuncMap) Get(method string) http.Handler {
	if handle, ok := rest[method]; ok {
		return handle
	}
	if handle, ok := rest[""]; ok {
		return handle
	}
	return nil
}

type regexHandleFunc struct {
	url   string
	regex *regexp.Regexp
	rest  restFuncMap
}

type RestHandler struct {
	normal map[string]restFuncMap
	regex  []*regexHandleFunc
}

// Get 注册GET方法处理器
func (rest *RestHandler) Get(url string, handler http.Handler) {
	rest.Add("get", url, handler)
}

// Post 注册POST方法处理器
func (rest *RestHandler) Post(url string, handler http.Handler) {
	rest.Add("post", url, handler)
}

// Put 注册PUT方法处理器
func (rest *RestHandler) Put(url string, handler http.Handler) {
	rest.Add("put", url, handler)
}

// Delete 注册DELETE方法处理器
func (rest *RestHandler) Delete(url string, handler http.Handler) {
	rest.Add("delete", url, handler)
}

// Add 注册GET方法处理器
// 如果url以`r:`开头，说明是一个正则表达式路径
func (rest *RestHandler) Add(method, url string, handlerFunc http.Handler) {
	if strings.HasPrefix(url, "r:") {
		url := strings.TrimPrefix(url, "r:")
		if rest.regex == nil {
			rest.regex = []*regexHandleFunc{}
		}
		var regexHandle *regexHandleFunc
		for _, regex := range rest.regex {
			if regex.url == url {
				regexHandle = regex
			}
		}
		if regexHandle != nil {
			regexHandle.rest[method] = handlerFunc
		} else {
			rest.regex = append(rest.regex, &regexHandleFunc{
				url:   url,
				regex: regexp.MustCompile(url),
				rest:  restFuncMap{method: handlerFunc}})
		}
	}
	if rest.normal == nil {
		rest.normal = map[string]restFuncMap{}
	}
	if restMap, ok := rest.normal[url]; ok {
		restMap[strings.ToLower(method)] = handlerFunc
	} else {
		rest.normal[url] = restFuncMap{method: handlerFunc}
	}
}

// ServeHTTP 实现Handle接口
func (rest *RestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.RequestURI
	method := strings.ToLower(r.Method)
	findButNotMatchMethod := false
	if restMap, ok := rest.normal[url]; ok {
		handle := restMap.Get(method)
		if handle != nil {
			handle.ServeHTTP(w, r)
			return
		} else {
			findButNotMatchMethod = true
		}
	}
	for _, regexFunc := range rest.regex {
		if regexFunc.regex.MatchString(url) {
			handle := regexFunc.rest.Get(method)
			if handle != nil {
				handle.ServeHTTP(w, r)
				return
			} else {
				findButNotMatchMethod = true
			}
		}
	}
	if findButNotMatchMethod {
		w.WriteHeader(405)
		w.Write([]byte("405 Method Not Support"))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 Not Found"))
	}
}

type HandlerFunc func(Context)

func createHandleFunc(handler HandlerFunc, medias ...media) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{w: w, r: r}
		for _, media := range medias {
			err := media(ctx)
			if err != nil {
				ctx.WriteString(err.Error())
				return
			}
		}
		handler(ctx)
	}
}

func createHandle(handler http.Handler, medias ...media) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{w: w, r: r}
		for _, media := range medias {
			err := media(ctx)
			if err != nil {
				ctx.WriteString(err.Error())
				return
			}
		}
		handler.ServeHTTP(w, r)
	}
}

type media func(context Context) error
