<html prefix="og: http://ogp.me/ns#">
<header>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<link rel="stylesheet" type="text/css" href="/static/css/main.css">

<link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
<link rel="manifest" href="/site.webmanifest">
<link rel="mask-icon" href="/safari-pinned-tab.svg" color="#5bbad5">
<meta name="msapplication-TileColor" content="#2b5797">
<meta name="theme-color" content="#ff6f00">

<meta property="og:site_name" content="BGS Log Viewer" />
<meta property="og:title" content="{{.Title}}" />
<meta property="og:type" content="{{.Type}}" />
<meta property="og:url" content="{{.Url}}" />
<meta property="og:description" content="{{.Description}}" />
{{if .HasImage}}<meta property="og:image" content="{{.ImageUrl}}" />{{end}}

<title>BGS Log Viewer</title>
</header>
<body>

  <h1>BGS Log Viewer</h1>

  <form action="./system", method="get" accept-charset="utf-8">
    <input type="text", name="q", maxlength="32" placeholder="System Name">
    <input type="submit">
  </form>

</body>
</html>
