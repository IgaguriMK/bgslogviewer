<html prefix="og: http://ogp.me/ns#">
<header>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<link rel="stylesheet" type="text/css" href="/static/css/system.css">

<link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
<link rel="manifest" href="/site.webmanifest">
<link rel="mask-icon" href="/safari-pinned-tab.svg" color="#5bbad5">
<meta name="msapplication-TileColor" content="#2b5797">
<meta name="theme-color" content="#ff6f00">

{{with .OGP}}
<meta property="og:site_name" content="BGS Log Viewer" />
<meta property="og:title" content="{{.Title}}" />
<meta property="og:type" content="{{.Type}}" />
<meta property="og:url" content="{{.Url}}" />
<meta property="og:description" content="{{.Description}}" />
{{if .HasImage}}<meta property="og:image" content="{{.ImageUrl}}" />{{end}}
{{end}}


<title>{{.SystemName}} -- BGS Log Viewer</title>
</header>
<body>

<a class="top" href="/">[TOP]</a>

<p>
Data source: <a href="https://www.edsm.net/">Elite Dangerous Star Map</a><br>
Cached at {{.CachedAt}}
</p>

<h1>{{.SystemName}}</h1>

<h2>Factions Overview</h2>
<table>
	<tr> <th colspan="4">Name</th> <th>Government</th> <th>Influence</th> <th>State</th> <th>Recovering</th> <th>Pending</th> <th>Last Update</th> </tr>
{{range .Overview}}
	<tr>
		<th class="con_right">{{if .IsControl}}<img width="16" height="16" src="/static/img/control.png">{{end}}</th>
		<th class="con_right con_left">{{if .IsPF}}<img width="16" height="16" src="/static/img/pf.png">{{end}}</th>
		<th class="con_right con_left">{{if .HasAllegiance}}<img width="16" height="16" src="{{.AllegianceImg}}">{{end}}</th>
		<th class="con_left"><a href="#{{.NameHash}}">{{.Name}}</a></th>
		<td>{{.Government}}</td>
		<td>{{.Influence}}</td>
		<td>{{.State}}</td>
		<td>{{.Recovering}}</td>
		<td>{{.Pending}}</td>
		<td>{{.LastUpdate}}</td>
	</tr>
{{end}}
</table>

{{if .RetreatedExists}}
<h3>Retreated</h3>

<ul>
{{range .Retreated}}
<li><a href="#{{.NameHash}}">{{.Name}}</a></li>
{{end}}
</ul>
{{end}}

<h2>Histories</h2>
{{range .History}}

<h3 id="{{.NameHash}}">{{.Name}}</h3>

<table class="stateHistory">
	<tr> <th>Date</th> <th>Influence</th> <th>State</th> <th>Recovering</th> <th>Pending</th> </tr>
	{{range .Records}}
	<tr>
		<th>{{.Date}}</th>
		{{with .Influence}}{{if .IsValid}}<td class="stateValid">{{.Value}}</td>{{else}}<td class="stateInvalid"></td>{{end}}{{end}}
		{{with .State}}{{if .IsValid}}<td class="stateValid">{{.Value}}</td>{{else}}<td class="stateInvalid"></td>{{end}}{{end}}
		{{with .Recovering}}{{if .IsValid}}<td class="stateValid">{{.Value}}</td>{{else}}<td class="stateInvalid"></td>{{end}}{{end}}
		{{with .Pending}}{{if .IsValid}}<td class="stateValid">{{.Value}}</td>{{else}}<td class="stateInvalid"></td>{{end}}{{end}}
	</tr>
	{{end}}
</table>

{{end}}

</body>
</html>
