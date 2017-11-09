package api

import (
	"config"
)

// 获取整体配置信息
func getConfig(ctx Context) {
	ctx.WriteJson(config.OriginConfig)
}

func getApiConfig(ctx Context) {
	ctx.WriteJson(config.OriginConfig.ApiInfo)
}

func getRulesConfig(ctx Context) {
	ctx.WriteJson(config.OriginConfig.Rules)
}

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
