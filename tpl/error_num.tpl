<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8">
        <title>错误码记录信息</title>
    </head>
    <body>
        <div>
            <p>总数: {{ .hits.total }}</p>
        </div>
        <div>
            <table border="1" cellspacing="0" cellpadding="0">
                <tr>
                    <td>doc id</td>
                    <td>用户</td>
                    <td>异常</td>
                    <td>错误信息</td>
                    <td>错误码</td>
                </tr>
                {{range .hits.hits}}
                <tr>
                    <td>{{ ._id }}</td>
                    <td>{{ ._source.message.token.id }}</td>
                    <td>{{ ._source.message.serviceException }}</td>
                    <td>{{ ._source.message.errorInfo }}</td>
                    <td>{{ ._source.message.code }}</td>
                </tr>
                {{else}}
                <tr>null</tr>
                {{end}}
            </table>
        </div>
    </body>
</html>