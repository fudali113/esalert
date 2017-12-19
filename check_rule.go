package main

import (
	"flag"
	"mylog"
	"encoding/json"
	"config"
	"rule"
)

func checkRule() {
	var ruleDir string
	flag.StringVar(&ruleDir, "ruleDir", "./sample/rules/2000_4000_code_range_count.yml", "配置文件所在目录")

	ruleConfig, err := config.GetRuleConfig(ruleDir)
	if err != nil {
		mylog.Error(err)
		return
	}
	testData := map[string]interface{}{}
	json.Unmarshal([]byte(ruleConfig.Test.Data), &testData)

	res, err := rule.NeedAlert(testData, ruleConfig.Script)
	if err != nil {
		mylog.Error("运行js脚本出错", err)
		return
	}
	if res != ruleConfig.Test.Should {
		mylog.Error("脚本运行结果与预期不一样")
		return
	}
	mylog.Info("测试通过")
}
