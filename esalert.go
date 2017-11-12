package main

import (
	"config"
	"rule"
	"mylog"
	"os"
	"flag"
	"api"
	"os/signal"
	"io/ioutil"
	"fmt"
	"util"
	"plugin"
	"storage"
	"alert"
)

var (
	configDir  string
	rulesName  string
	configName string
	pluginDir  string
)

func main() {
	flag.StringVar(&configDir, "config", "./sample", "配置文件所在目录")
	flag.StringVar(&rulesName, "rules_name", "rules", "rule配置文件夹的名字")
	flag.StringVar(&configName, "config_name", "config.yml", "配置文件的名字")
	flag.StringVar(&pluginDir, "plugin", "./plugin", "插件文件所在目录")
	config, err := config.IntiConfig(config.ConfigDirInfo{Dir: configDir, RuleName: rulesName, ConfigName: configName})
	checkErr(err)
	err = rule.Run(*config)
	checkErr(err)
	//checkErr(loadPlugin(pluginDir))
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

func loadPlugin(dir string) error {
	pluginDir, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("打开插件文件夹"+dir+"夹出错, err : %s", err.Error())
	}
	for _, pluginFile := range pluginDir {
		if !pluginFile.IsDir() {
			pluginFileName := util.BuildFileDir(dir, pluginFile.Name())
			plugin, err := plugin.Open(pluginFileName)
			if err != nil {
				return fmt.Errorf("加载插件"+pluginFileName+"出错, err : %s", err.Error())
			}
			storagePlugin, err := plugin.Lookup("Storage")
			if err == nil {
				if sp, ok := storagePlugin.(storage.StorageCreater); ok {
					storage.Register(sp)
				}
			}
			alertPlugin, err := plugin.Lookup("Alert")
			if err == nil {
				if ap, ok := alertPlugin.(alert.AlerterCreater); ok {
					alert.Register(ap)
				}
			}
		}
	}
	return nil
}
