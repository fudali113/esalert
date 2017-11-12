package api

import (
	"rule"
	"io/ioutil"
	"encoding/json"
	"config"
)

// 获取context相关信息
func getContext(ctx Context) {
	ctx.WriteJson(rule.RuleContextMap)
}

// 新增一条rule
func postRule(ctx Context) {
	bytes, err := ioutil.ReadAll(ctx.r.Body)
	if err != nil {
		ctx.WriteError(err)
		return
	}
	ruleConfig := config.RuleConfig{}
	json.Unmarshal(bytes, &ruleConfig)
	_, err = rule.RunRule(ruleConfig)
	if err != nil {
		ctx.WriteError(err)
		return
	}
	ctx.WriteString("success")
}

func ruleStart(ctx Context) {
	ruleName := ctx.URISubSelect(-2)
	if rc, ok := rule.RuleContextMap[ruleName]; ok {
		if rc.Status == rule.START {
			ctx.WriteString("改rule已经启动，如想重启请访问restart接口")
			return
		}
		rc.Start()
		return
	}
	ctx.Four04()
}

func ruleStop(ctx Context) {
	ruleName := ctx.URISubSelect(-2)
	if rc, ok := rule.RuleContextMap[ruleName]; ok {
		rc.Stop()
		return
	}
	ctx.Four04()
}

func ruleRestart(ctx Context) {
	ruleName := ctx.URISubSelect(-2)
	if rc, ok := rule.RuleContextMap[ruleName]; ok {
		rc.Restart()
		return
	}
	ctx.Four04()
}
