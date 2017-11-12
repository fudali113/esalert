<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8">
        <title>stack异常提醒</title>
    </head>
    <body>
        <div>
            <p>此报警没一分钟运行一次，查询前两分钟的日志信息(以为日志同步有延迟，所以设置两分钟，所有您有可能会收到两份一样的报警，但是林可错杀一千，我们也绝不错过，注意查看文档的id信息确定是否为同一份报警)</p>
            <p>总数: {{ .hits.total }}</p>
        </div>
        <div>
            <table>
            {{range .hits.hits}}
            <tr>
            <p>文档id: {{ ._id }}</p>
            <p>错误码: {{ ._source.message.code }}</p>
            <p>用户id: {{ ._source.message.token.id }}</p>
            <p>请求路径: {{ ._source.message.errorInfo.path }}</p>
            <p>请求参数: {{ ._source.message.params }}</p>
            <p>错误信息: {{ ._source.message.errorInfo.message }}</p>
            <p>异常信息: {{ ._source.message.errorInfo.error }}</p>
            <p>异常堆栈: {{ ._source.message.stack }}</p>
            </tr>
            {{end}}
            </table>
        </div>
    </body>
</html>