package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/IgaguriMK/bgslogviewer/apiCaller"
	"github.com/IgaguriMK/bgslogviewer/config"
	"github.com/IgaguriMK/bgslogviewer/prof"
	"github.com/IgaguriMK/bgslogviewer/view"
)

const (
	systemNameLimit = 64
)

const (
	cacheLocalMin = 120 * time.Second
	cacheLocalMax = 60 * time.Minute
	cacheCdnMin   = 5 * time.Minute
	cacheCdnMax   = 6 * time.Hour
)

func main() {
	time.Sleep(3 * time.Second)

	r := gin.Default()

	r.StaticFile("/", "./static/main.html")
	r.StaticFile("/index.html", "./static/main.html")

	r.GET("/system", statPage)

	r.Static("/static/css", "./static/css")
	r.Static("/static/img", "./static/img")

	r.StaticFile("/android-chrome-192x192.png", "./static/favicon/android-chrome-192x192.png")
	r.StaticFile("/android-chrome-512x512.png", "./static/favicon/android-chrome-512x512.png")
	r.StaticFile("/apple-touch-icon.png", "./static/favicon/apple-touch-icon.png")
	r.StaticFile("/browserconfig.xml", "./static/favicon/browserconfig.xml")
	r.StaticFile("/favicon.ico", "./static/favicon/favicon.ico")
	r.StaticFile("/favicon-16x16.png", "./static/favicon/favicon-16x16.png")
	r.StaticFile("/favicon-32x32.png", "./static/favicon/favicon-32x32.png")
	r.StaticFile("/mstile-150x150.png", "./static/favicon/mstile-150x150.png")
	r.StaticFile("/safari-pinned-tab.svg", "./static/favicon/safari-pinned-tab.svg")
	r.StaticFile("/site.webmanifest", "./static/favicon/site.webmanifest")

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("[FATAL] Can't execute gin server: ", err)
	}
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

	if systemName == "" || len(systemName) > systemNameLimit {
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

	switch stat {
	case apiCaller.Timeout:
		c.String(503, "Too many access")
		pf.End(503)
		return
	case apiCaller.Invalid:
		c.String(404, "Invalid request")
		pf.End(404)
		return
	}

	if systemName != v.Name {
		params := url.Values{}
		params.Add("q", v.Name)
		u := "/system?" + params.Encode()

		c.Redirect(301, u)
		return
	}

	pf.Mark("template")
	body := new(bytes.Buffer)

	err = view.System(body, v)
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
