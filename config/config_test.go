package config

import (
	"testing"
	"mylog"
	"util"
)

func Test_IntiConfig(t *testing.T) {
	config, err := IntiConfig(ConfigDirInfo{Dir: "../sample", RuleName: "rules", ConfigName: "config.yml"})
	if err != nil {
		t.Error(err)
	}
	if len(config.Rules) == 0 {
		t.Error("解析出错")
	}
	for _, rule := range config.Rules {
		if rule.Storage == nil {
			t.Error("解析出错")
		}
		json, err := util.QueryToJSON(rule.Storage)
		mylog.Info(string(json), err)
	}
}
