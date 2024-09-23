package consts

const ErrorEmail = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
        }
        .info-container {
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background-color: #f9f9f9;
            border: 1px solid #ddd;
            border-radius: 5px;
        }
        .info-item {
            margin-bottom: 10px;
        }
        .content {
            white-space: pre-line;
        }
    </style>
</head>
<body>
<div class="info-container">
    <div class="info-item"><strong>时间：</strong>{GEN_TIME}</div>
    <div class="info-item">
        <strong>内容：</strong>
        <code class="content">
            {SOME_MESSAGE}
        </code>
    </div>
</div>
</body>
</html>
`
