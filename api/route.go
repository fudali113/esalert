package api

import (
	"config"
	"fmt"
	"net/http"
	"mylog"
)

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
	handler.Get("/api/config", createHandleFunc(getConfig, basicAuthMedia))
	handler.Get("/api/config/rules", createHandleFunc(getRulesConfig, basicAuthMedia))
	handler.Get("r:/api/config/rules/[\\S\\s]+", createHandleFunc(getOneRulesConfig, basicAuthMedia))
	handler.Get("/api/config/api-info", createHandleFunc(getApiConfig, basicAuthMedia))

	handler.Get("/api/context", createHandleFunc(getContext, basicAuthMedia))
	handler.Post("/api/context", createHandleFunc(postRule, basicAuthMedia))

	handler.Post("r:/api/context/[\\S\\s]+/stop", createHandleFunc(ruleStop, basicAuthMedia))
	handler.Post("r:/api/context/[\\S\\s]+/start", createHandleFunc(ruleStart, basicAuthMedia))
	handler.Post("r:/api/context/[\\S\\s]+/restart", createHandleFunc(ruleRestart, basicAuthMedia))
	// 静态文件处理
	handler.Get("r:/[\\s\\S]*", createHandle(http.FileServer(http.Dir("./static/")), basicAuthMedia))

	go func() {
		err := http.ListenAndServe(":"+config.Port, handler)
		if err != nil {
			mylog.Error("启动api服务失败", err)
		}
	}()
	mylog.Info("启动api服务中", config.Port)
}
