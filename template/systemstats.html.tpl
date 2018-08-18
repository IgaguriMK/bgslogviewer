<html>
<header>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<title>{{.Name}} -- BGS Log Viewer</title>
</header>
<body>

<a class="top" href="/">[TOP]</a>

<h1>{{.Name}}</h1>

<h2>概況</h2>
<table>
	<tr> <th>名前</th> <th>Influence</th> <th>State</th> <th>Recovering</th> <th>Pending</th> <th>最終更新</th> </tr>
{{range $i, $f := .Factions}}
	<tr class=".{{oddEven $i}}">
		<td>{{$f.Name}}</td>
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

<h2>履歴</h2>
{{range .Factions}}

<h3>{{.Name}}</h3>

<table>
	<tr> <th>日時</th> <th>Influence</th> <th>State</th> <th>Recovering</th> <th>Pending</th> </tr>
{{range $i, $h := .History}}
	<tr class=".{{oddEven $i}}">
		<td>{{$h.DateStr}}</td>
		<td>{{$h.InfluenceStr}}</td>
		<td>{{$h.Current}}</td>
		<td>{{range $h.Recovering}}{{.State}}{{.TrendStr}} {{end}} </td>
		<td>{{range $h.Pending}}{{.State}}{{.TrendStr}} {{end}} </td>
	</tr>
{{end}}
</table>

{{end}}

</body>
</html>
