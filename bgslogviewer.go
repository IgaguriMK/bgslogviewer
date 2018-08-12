package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"time"

	"github.com/IgaguriMK/bgslogviewer/apiCaller"
	"github.com/gin-gonic/gin"
)

const (
	timeFormat = "2006-01-02 (15)"
)

var mainTemplate *template.Template
var statTemplate *template.Template

func init() {
	var err error

	statTemplate, err = template.ParseFiles("template/systemstats.html.tpl")
	if err != nil {
		log.Fatal("[FATAL] Failed parse template: ", err)
	}
}

func main() {
	r := gin.Default()

	r.GET("/", mainPage)
	r.GET("/system", statPage)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("[FATAL] Can't execute gin server: ", err)
	}
}

func mainPage(c *gin.Context) {
	c.File("static/main.html")
}

func statPage(c *gin.Context) {
	systemName := c.Query("q")

	if systemName == "" {
		c.String(404, "Invalid query")
		return
	}

	v, stat, err := apiCaller.FetchFactions(systemName)
	if err != nil {
		c.String(500, "Internal error")
		log.Println("[ERROR] fetching data error: ", err)
		return
	}

	if stat == apiCaller.Invalid {
		c.String(404, "Invalid request")
		return
	}

	v.GenStr(time.UTC, timeFormat)

	body := new(bytes.Buffer)

	err = statTemplate.Execute(body, v)
	if err != nil {
		c.String(500, "Internal error")
		log.Fatal("[FATAL] Execute template: ", err)
		return
	}

	c.Status(200)

	c.Header("Cache-Control", "max-age=600, s-maxage=600")

	_, err = io.Copy(c.Writer, body)
	if err != nil {
		log.Println("[INFO] error while sending data: ", err)
	}
}
