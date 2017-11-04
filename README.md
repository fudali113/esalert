# esalert                   
[![Build Status](https://travis-ci.org/23mf/esalert.svg?branch=master)](https://travis-ci.org/23mf/esalert)
[![Go Report Card](https://goreportcard.com/badge/github.com/23mf/esalert)](https://goreportcard.com/report/github.com/23mf/esalert)

提供查询elasticsearch数据根据规则报警功能

查询语句自定义，可以自己编写es查询query，只需要会使用es即可，没有其他任何学习成本
是否报警规则使用js脚本语言判断，只需要会使用js即可

简单明了，工具本身学习成本很低，只需要会使用es与js即可简单上手

# 超级简单的配置
```yaml
host: localhost      # es host
port: 9200          # es port
username: elastic   # es username
password: changeme  # es password
mail:
  username: fuyi@23mofang.com
  password: xxxxxxxxxx
  smtp_host: smtp.exmail.qq.com
  smtp_port: 25
  send_to:
    - fuyi@23mofang.com
#    - yuxiaobo@23mofang.com
  from_addr: fuyi@23mofang.com             # 显示发送出去的用户是谁
  reply_to: fuyi@23mofang.com              # 发送出去的邮件回复给谁
#  tpl_file: "/xx/xx/xx.tpl"         # go template模板文件     tpl_file与content必须存在一个
#  content: "xxx{{total}}xxxx"       # go template模板字符串
#  theme: "xxxx"                     # 邮件主题
rules:              # 检查规则
  - name: test   # 没有规则必须有一个唯一的name
    index: gateway-*
    body:
      query:
        bool:
          must:
          - exists:
              field: message.serviceException
          - range:
              "@timestamp":
                gte: now-30m
    script: '''
            res.hits.total > 10
            '''        # 默认会将查询获取的json数据易以res变量在脚本作用域内， 当该脚本返回true时执行报警
    interval:       # 隔多久发起一次请求，该字段会根据里面的语义信息转换时间
      minute: 30
    alerts:                                  # 报警
#      - type: http                          # http报警规则
#        url: http://baidu.com
      - type: mail                          # mail报警规则
        mail:                               # 该配置项参数与外层mail参数一致，该配置优先级高于外层mail配置
          tpl_file: error_num.tpl         # go template模板文件     tpl_file与content必须存在一个
#          content: "xxx{{total}}xxxx"       # go template模板字符串
          subject: esalert test                     # 邮件主题
  - name: error_code
    index: gateway-*
    body:
      size: 0
      query:
        range:
          "@timestamp":
            gte: now-30m
      aggs:
        price_ranges:
          range:
            field: message.code
            ranges:
            - to: 1000
            - from: 1000
              to: 2000
            - from: 2000
              to: 4000
            - from: 4000
              to: 5000
            - from: 5000
              to: 6000
            - from: 6000
    script: true
    interval:
      minute: 30
    alerts:
      - type: mail
        mail:
          tpl_file: tpl/error_code.tpl
          subject: error code agg
```

# Futures
* 更加完善的日志记录
* 提供能多报警方式
* 使每个运行的rule可管理并可灵活扩充
* 提供web界面
