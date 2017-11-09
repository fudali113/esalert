<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8">
        <title>stack异常提醒</title>
    </head>
    <body>
        <div>
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