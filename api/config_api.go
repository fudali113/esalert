package api

import (
	"config"
)

// 获取原始配置信息
func getConfig(ctx Context) {
	ctx.WriteJson(config.OriginConfig)
}

// 获取原始api配置信息
func getApiConfig(ctx Context) {
	ctx.WriteJson(config.OriginConfig.ApiInfo)
}

// 获取原始rule配置信息
func getRulesConfig(ctx Context) {
	ctx.WriteJson(config.OriginConfig.Rules)
}

// 获取某条rule原始配置信息
func getOneRulesConfig(ctx Context) {
	name := ctx.URISubSelect(-1)
	for _, rule := range config.OriginConfig.Rules {
		if rule.Name == name {
			ctx.WriteJson(rule)
			return
		}
	}
	ctx.w.Write([]byte("name为" + name + "的rule不存在"))
	ctx.w.WriteHeader(404)
}
