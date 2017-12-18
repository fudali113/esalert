# esalert                   
[![Build Status](https://travis-ci.org/23mf/esalert.svg?branch=master)](https://travis-ci.org/23mf/esalert)
[![Go Report Card](https://goreportcard.com/badge/github.com/23mf/esalert)](https://goreportcard.com/report/github.com/23mf/esalert)

提供查询elasticsearch数据根据规则报警功能

查询语句自定义，可以自己编写es查询query，只需要会使用es即可，没有其他任何学习成本
是否报警规则使用js脚本语言判断，只需要会使用js即可

简单明了，工具本身学习成本很低，只需要会使用es与js即可简单上手

# 超级简单的配置
```yaml
storage:
  _type: es
  host: localhost     # es host
  port: 9200          # es port
  username: elastic   # es username
  password: changeme  # es password
  index: gateway-*
api:
  enable: true
  port: 3131
  basic_auth:
    enable: true
    username: admin
    password: 123456
alert:
  _type: mail             # alert type
  username: fudali4test@163.com
  password: 1234567890abc
  smtp_host: smtp.163.com
  smtp_port: 25
  send_to:
    - fudali4test@163.com
  from_addr: fudali4test@163.com             # this email from who
  reply_to: fudali4test@163.com              # this email rrply to who
rules:              # rule policy 
  - name: exists_stack_alert   # rule name , must unique
    storage: 
      index: gateway-*
      body:
        query:
          bool:
            must:
            - exists:
                field: message.stack
            - range:
                "@timestamp":
                  gte: now-2m
    # 默认会将查询获取的json数据易以`result`变量在脚本作用域内， 当该脚本返回true时执行报警
    script: >
            result.hits.total > 0
    interval:       # 隔多久发起一次请求，该字段会根据里面的语义信息转换时间
      m: 1
    alerts:                                  # 报警
    #      - type: http                      # http报警规则
    #        url: http://baidu.com
      - tpl_file: sample/tpl/exists_stack_alert.tpl         # go template模板文件     tpl_file与content必须存在一个
        content: "xxx{{total}}xxxx"                 # go template模板字符串
        subject: 错误异常堆栈提醒                     # 邮件主题
        send_to:
          - fuyi@23mofang.com

```

# Futures
* optimize log
* more alerter
* rule manager
* support web ui

# build
```
git clone github.com/fudali113/esalert esalert/src
build: 
GOPATH=./esalert ./build.sh
run:
GOPATH=./esalert go run esalert.go
```
