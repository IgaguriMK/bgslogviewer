package main

import (
	"log"
	"time"

	"github.com/IgaguriMK/bgslogviewer/controller"
	"github.com/gin-gonic/gin"
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

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("[FATAL] Can't execute gin server: ", err)
	}
}
