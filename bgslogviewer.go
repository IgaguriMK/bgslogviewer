package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/IgaguriMK/bgslogviewer/apiCaller"
	"github.com/IgaguriMK/bgslogviewer/config"
	"github.com/IgaguriMK/bgslogviewer/prof"
)

const (
	timeFormat = "2006-01-02 (15)"
)

const (
	cacheLocalMin = 120 * time.Second
	cacheLocalMax = 60 * time.Minute
	cacheCdnMin   = 5 * time.Minute
	cacheCdnMax   = 6 * time.Hour
)

var mainTemplate *template.Template
var statTemplate *template.Template

func init() {
	var err error

	funcs := template.FuncMap{
		"oddEven": func(n int) string {
			if (n+1)%2 == 0 {
				return "even"
			}
			return "odd"
		},
	}

	statTemplate, err = template.New("systemstats.html.tpl").Funcs(funcs).ParseFiles("template/systemstats.html.tpl")
	if err != nil {
		log.Fatal("[FATAL] Failed parse template: ", err)
	}
}

func main() {
	r := gin.Default()

	r.GET("/", mainPage)
	r.GET("/index.html", mainPage)
	r.GET("/system", statPage)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("[FATAL] Can't execute gin server: ", err)
	}
}

func mainPage(c *gin.Context) {
	commonHeader(c)

	c.Status(200)
	c.File("static/main.html")
}

func statPage(c *gin.Context) {
	pf := prof.NewProfiler()
	pf.Start("validation")

	systemName, ok := c.GetQuery("q")
	if !ok {
		c.String(400, "Invalid Url")
		return
	}

	pf.AddParam("q", systemName)

	commonHeader(c)

	if systemName == "" {
		c.String(404, "Invalid query")
		pf.End(404)
		return
	}

	pf.Mark("fetchAPI")
	v, stat, err := apiCaller.FetchFactions(systemName)
	if err != nil {
		c.String(500, "Internal error")
		pf.End(500)
		log.Println("[ERROR] fetching data error: ", err)
		return
	}

	if stat == apiCaller.Invalid {
		c.String(404, "Invalid request")
		pf.End(404)
		return
	}

	pf.Mark("genStr")
	v.GenStr(time.UTC, timeFormat)

	pf.Mark("template")
	body := new(bytes.Buffer)

	err = statTemplate.Execute(body, v)
	if err != nil {
		c.String(500, "Internal error")
		pf.End(500)
		log.Fatal("[FATAL] Execute template: ", err)
		return
	}

	c.Status(200)

	ll, cl := cachesLen()
	c.Header("Cache-Control", fmt.Sprintf("max-age=%d, s-maxage=%d", ll, cl))

	pf.Mark("send")
	_, err = io.Copy(c.Writer, body)
	if err != nil {
		pf.End(-1)
		log.Println("[INFO] error while sending data: ", err)
	}

	pf.End(200)
}

func commonHeader(c *gin.Context) {
	c.Header("X-XSS=Protection", "1; mode=block")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Security-Policy", "default-src 'self'")
}

func cachesLen() (localCache int, cdnCache int) {
	now := time.Now()

	nextUpdate := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		config.BgsUpdate,
		0,
		0,
		0,
		time.UTC,
	)

	if nextUpdate.Before(now) {
		nextUpdate = nextUpdate.Add(24 * time.Hour)
	}

	prevUpdate := nextUpdate.Add(-24 * time.Hour)

	// 直前のBGS更新から近いときは最短キャッシュ
	if now.Add(-3 * time.Hour).Before(prevUpdate) {
		return int(cacheLocalMin.Seconds()), int(cacheCdnMin.Seconds())
	}

	// 次のBGS更新まで間近
	if now.Add(time.Hour).After(nextUpdate) {
		return int(cacheLocalMin.Seconds()), int(cacheCdnMin.Seconds())
	}

	// BGS更新付近ではないが、最大キャッシュでは長すぎる
	if now.Add(cacheCdnMax + cacheCdnMin + time.Hour).After(nextUpdate) {
		d := nextUpdate.Sub(now) - time.Hour
		return int(d.Seconds()), int(d.Seconds())
	}

	// 次のBGS更新まで余裕がある
	return int(cacheLocalMax.Seconds()), int(cacheCdnMax.Seconds())
}
