package main

import (
	"config"
	"rule"
	"mylog"
	"os"
	"flag"
	"api"
	"os/signal"
)

var (
	configDir  string
	rulesName  string
	configName string
)

func main() {
	flag.StringVar(&configDir, "sample", "./sample", "配置文件所在目录")
	flag.StringVar(&rulesName, "rules_name", "rules", "rule配置文件夹的名字")
	flag.StringVar(&configName, "config_name", "config.yml", "配置文件的名字")
	config, err := config.IntiConfig(config.ConfigDirInfo{Dir: configDir, RuleName: rulesName, ConfigName: configName})
	checkErr(err)
	err = rule.Run(*config)
	checkErr(err)
	// 运行api服务
	api.Start(config.ApiInfo)
	mylog.Info("start success")
	// 保证主进程在接收到退出信号前不退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	mylog.Info("Got signal:", s)
}

func checkErr(err error) {
	if err != nil {
		mylog.Error("解析配置文件出错, ", err)
		os.Exit(1)
	}
}
