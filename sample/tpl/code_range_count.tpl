<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8">
        <title>错误码聚合骤然上升提醒</title>
    </head>
    <body>
        <div>
            <table border="1" cellspacing="0" cellpadding="0">
                <tr>
                    <td>时间</td>
                    <td>总数</td>
                </tr>
                {{range .aggregations.code_count.buckets}}
                <tr>
                    <td>{{ .key_as_string }}</td>
                    <td>{{ .doc_count }}</td>
                </tr>
                {{else}}
                <tr>null</tr>
                {{end}}
            </table>
        </div>
    </body>
</html>