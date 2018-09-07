package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/IgaguriMK/bgslogviewer/controller"
	"github.com/gin-gonic/gin"
)

var logCh chan AccessLog

func main() {
	startLogSaver()

	time.Sleep(3 * time.Second)

	e := gin.Default()

	r := e.Use(accessLogger)

	r.StaticFile("/", "./static/main.html")
	r.StaticFile("/index.html", "./static/main.html")

	r.GET("/system", controller.StatPage)

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

	r.StaticFile("/robots.txt", "./static/misc/robots.txt")

	err := e.Run(":8080")
	if err != nil {
		log.Fatal("[FATAL] Can't execute gin server: ", err)
	}
}

type AccessLog struct {
	Date      string `json:"date"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Code      int64  `json:"code"`
	From      string `json:"from"`
	UserAgent string `json:"useragent"`
	Duration  int64  `json:"dur_ms"`
}

func accessLogger(c *gin.Context) {
	start := time.Now()

	c.Next()

	dur := time.Since(start)

	logCh <- AccessLog{
		Date:      start.Format(time.RFC3339),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Code:      int64(c.Writer.Status()),
		From:      c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Duration:  int64(dur / time.Millisecond),
	}
}

func startLogSaver() {
	err := os.MkdirAll("./log", 0744)
	if err != nil {
		log.Fatal("alert: can't create log direcrory: ", err)
	}

	logCh = make(chan AccessLog, 8)
	go func() {
		for l := range logCh {
			saveLog(l)
		}
	}()
}

func saveLog(l AccessLog) {
	f, err := os.OpenFile("./log/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("error: can't open access.log: ", err)
		return
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(l)
	if err != nil {
		log.Println("error: access log save error:", err)
	}
}
