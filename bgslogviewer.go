package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/comail/colog"
	"github.com/gin-gonic/gin"

	"github.com/IgaguriMK/bgslogviewer/controller"
)

func main() {
	startLogSaver()

	time.Sleep(3 * time.Second)

	e := gin.Default()

	r := e.Use(
		controller.CommonHeader,
	)

	r.GET("/", controller.MainPage)
	r.GET("/system", controller.StatPage)

	rs := r.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "max-age=600, s-maxage=604800")
		c.Next()
	})

	rs.Static("/static/css", "./static/css")
	rs.Static("/static/img", "./static/img")

	rs.StaticFile("/android-chrome-192x192.png", "./static/favicon/android-chrome-192x192.png")
	rs.StaticFile("/android-chrome-512x512.png", "./static/favicon/android-chrome-512x512.png")
	rs.StaticFile("/apple-touch-icon.png", "./static/favicon/apple-touch-icon.png")
	rs.StaticFile("/browserconfig.xml", "./static/favicon/browserconfig.xml")
	rs.StaticFile("/favicon.ico", "./static/favicon/favicon.ico")
	rs.StaticFile("/favicon-16x16.png", "./static/favicon/favicon-16x16.png")
	rs.StaticFile("/favicon-32x32.png", "./static/favicon/favicon-32x32.png")
	rs.StaticFile("/mstile-150x150.png", "./static/favicon/mstile-150x150.png")
	rs.StaticFile("/safari-pinned-tab.svg", "./static/favicon/safari-pinned-tab.svg")
	rs.StaticFile("/site.webmanifest", "./static/favicon/site.webmanifest")

	rs.StaticFile("/robots.txt", "./static/misc/robots.txt")

	err := e.Run(":8080")
	if err != nil {
		log.Fatal("[FATAL] Can't execute gin server: ", err)
	}
}

func startLogSaver() {
	err := os.MkdirAll("./log", 0744)
	if err != nil {
		log.Fatal("alert: can't create log direcrory: ", err)
	}

	// error.log
	colog.Register()

	logf, err := os.OpenFile("./log/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("alert: can't open error.log: ", err)
	}
	logw := io.MultiWriter(logf, os.Stderr)
	colog.SetOutput(logw)
	if level, ok := os.LookupEnv("LOGLEVEL"); ok {
		lvl, err := colog.ParseLevel(level)
		if err != nil {
			log.Fatalf("alert: invalid LOGLEVEL = %q", level)
		}
		colog.SetMinLevel(lvl)
	} else {
		colog.SetMinLevel(colog.LInfo)
	}

	// access.log
	noLogPatterns := make([]*regexp.Regexp, 0)

	if f, err := os.Open("./conf.d/nolog-agent.txt"); err != nil {
		log.Println("info: './conf.d/nolog-agent.txt' not found. use empty list.")
	} else {
		sc := bufio.NewScanner(f)

		for sc.Scan() {
			pat := sc.Text()
			exp, err := regexp.Compile(pat)
			if err != nil {
				log.Fatalf("alert: invalid agent regexp %q: %v", pat, err)
			}
			noLogPatterns = append(noLogPatterns, exp)
			log.Printf("debug: add ignored useragent %s", pat)
		}

		f.Close()
	}
}

func matchPatterns(str string, exps []*regexp.Regexp) bool {
	for _, exp := range exps {
		if exp.MatchString(str) {
			return true
		}
	}
	return false
}
