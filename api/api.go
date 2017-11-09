package api

import (
	"net/http"
	"mylog"
	"fmt"
	"strings"
	"regexp"
	"config"
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
		w.Write([]byte("405 Method Not Support"))
		w.WriteHeader(405)
	} else {
		w.Write([]byte("404 Not Found"))
		w.WriteHeader(404)
	}
}

// Start 根据配置参数启动Api服务
func Start(config config.ApiConfig) {
	if !config.Enable {
		return
	}
	auth := config.BasicAuth
	var basicAuthMedia = func(ctx Context) error {
		if !auth.Enable {
			return nil
		}
		if name, pass, ok := ctx.r.BasicAuth(); ok && name == auth.Username && pass == auth.Password {
			return nil
		}
		ctx.w.Header().Add("WWW-Authenticate", `Basic realm="Login Required"`)
		ctx.w.WriteHeader(401)
		return fmt.Errorf("用户名密码错误")
	}
	handler := &RestHandler{}
	handler.Get("/api/sample", createHandleFunc(getConfig, basicAuthMedia))
	handler.Get("/api/sample/rules", createHandleFunc(getRulesConfig, basicAuthMedia))
	handler.Get("r:/api/sample/rules/[\\S\\s]+", createHandleFunc(getOneRulesConfig, basicAuthMedia))
	handler.Get("/api/sample/api-info", createHandleFunc(getApiConfig, basicAuthMedia))

	handler.Get("/api/context", createHandleFunc(getContext, basicAuthMedia))
	go func() {
		err := http.ListenAndServe(":"+config.Port, handler)
		if err != nil {
			mylog.Error("启动api服务失败", err)
		}
	}()
	mylog.Info("启动api服务中", config.Port)
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
