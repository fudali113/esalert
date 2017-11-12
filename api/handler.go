package api

import (
	"net/http"
	"strings"
	"regexp"
)

var defaultMedias = []media{}

type restFuncMap map[string]http.HandlerFunc

func (rest restFuncMap) Get(method string) http.HandlerFunc {
	if handleFunc, ok := rest[method]; ok {
		return handleFunc
	}
	if handleFunc, ok := rest[""]; ok {
		return handleFunc
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
func (rest *RestHandler) Get(url string, handlerFunc http.HandlerFunc) {
	rest.Add("get", url, handlerFunc)
}

// Post 注册POST方法处理器
func (rest *RestHandler) Post(url string, handlerFunc http.HandlerFunc) {
	rest.Add("post", url, handlerFunc)
}

// Put 注册PUT方法处理器
func (rest *RestHandler) Put(url string, handlerFunc http.HandlerFunc) {
	rest.Add("put", url, handlerFunc)
}

// Delete 注册DELETE方法处理器
func (rest *RestHandler) Delete(url string, handlerFunc http.HandlerFunc) {
	rest.Add("delete", url, handlerFunc)
}

// Add 注册GET方法处理器
// 如果url以`r:`开头，说明是一个正则表达式路径
func (rest *RestHandler) Add(method, url string, handlerFunc http.HandlerFunc) {
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
		handleFunc := restMap.Get(method)
		if handleFunc != nil {
			handleFunc(w, r)
			return
		} else {
			findButNotMatchMethod = true
		}
	}
	for _, regexFunc := range rest.regex {
		if regexFunc.regex.MatchString(url) {
			handleFunc := regexFunc.rest.Get(method)
			if handleFunc != nil {
				handleFunc(w, r)
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

type media func(context Context) error
