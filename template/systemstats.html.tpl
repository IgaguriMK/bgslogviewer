<html>
<header>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<link rel="stylesheet" type="text/css" href="/static/css/system.css">
<title>{{.Name}} -- BGS Log Viewer</title>
</header>
<body>

<a class="top" href="/">[TOP]</a>

<h1>{{.Name}}</h1>

<h2>Factions Overview</h2>
<table>
	<tr> <th>Name</th> <th>Influence</th> <th>State</th> <th>Recovering</th> <th>Pending</th> <th>Last Update</th> </tr>
{{range $i, $f := .Factions}}
	<tr class=".{{oddEven $i}}">
		<th>{{$f.Name}}</th>
		{{with $f.NewestStates}}
			<td>{{.InfluenceStr}}</td>
			<td>{{.Current}}</td>
			<td>{{range .Recovering}}{{.State}}{{.TrendStr}} {{end}} </td>
			<td>{{range .Pending}}{{.State}}{{.TrendStr}} {{end}} </td>
			<td>{{.DateStr}}</td>
		{{end}}
	</tr>
{{end}}
</table>

<h2>Histories</h2>
{{range .Factions}}

<h3>{{.Name}}</h3>

<table class="stateHistory">
	<tr> <th>Date</th> <th>Influence</th> <th>State</th> <th>Recovering</th> <th>Pending</th> </tr>
{{range $i, $h := .History}}
	<tr class=".{{oddEven $i}}">
		<th>{{$h.DateStr}}</th>
		<td class="{{if $h.ValidInfluence}}stateValid{{else}}stateInvalid{{end}}">{{$h.InfluenceStr}}</td>
		<td class="{{if $h.ValidCurrent}}stateValid{{else}}stateInvalid{{end}}">{{$h.Current}}</td>
		<td class="{{if $h.ValidRecovering}}stateValid{{else}}stateInvalid{{end}}">{{range $h.Recovering}}{{.State}}{{.TrendStr}}{{end}} </td>
		<td class="{{if $h.ValidPending}}stateValid{{else}}stateInvalid{{end}}">{{range $h.Pending}}{{.State}}{{.TrendStr}}{{end}} </td>
	</tr>
{{end}}
</table>

{{end}}

</body>
</html>
