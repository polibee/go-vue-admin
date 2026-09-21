package openapi

const scalarHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Go Vue Admin API</title>
  </head>
  <body>
    <script id="api-reference" data-url="/api/openapi.json" src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`

func DocsHTML() string { return scalarHTML }
