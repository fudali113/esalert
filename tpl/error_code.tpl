<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8">
        <title>错误码聚合</title>
    </head>
    <body>
        <div>
            <table border="1" cellspacing="0" cellpadding="0">
                <tr>
                    <td>key</td>
                    <td>from</td>
                    <td>to</td>
                    <td>count</td>
                </tr>
                {{range .aggregations.price_ranges.buckets}}
                <tr>
                    <td>{{ .key }}</td>
                    <td>{{ .from }}</td>
                    <td>{{ .to }}</td>
                    <td>{{ .doc_count }}</td>
                </tr>
                {{else}}
                <tr>null</tr>
                {{end}}
            </table>
        </div>
    </body>
</html>