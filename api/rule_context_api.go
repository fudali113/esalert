package api

import (
	"rule"
)

func getContext(ctx Context) {
	ctx.WriteJson(rule.RuleContextMap)
}
