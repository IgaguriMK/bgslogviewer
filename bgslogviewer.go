package main

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/IgaguriMK/bgslogviewer/api"
	"github.com/IgaguriMK/bgslogviewer/model"
)

const (
	timeFormat = "2006-01-02 (15)"
)

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

	r.GET("/", statPage)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("[FATAL] Can't execute gin server: ", err)
	}
}

func statPage(c *gin.Context) {
	f, err := os.Open("sol.json")
	if err != nil {
		c.String(500, "External API error")
		log.Println("[ERROR] loading data error: ", err)
		return
	}
	defer f.Close()

	var res api.SystemFactions
	err = json.NewDecoder(f).Decode(&res)
	if err != nil {
		c.String(500, "External API error")
		log.Println("[ERROR] loading data error: ", err)
		return
	}

	v := model.FromApiResult(res)
	v.GenStr(time.UTC, timeFormat)

	err = statTemplate.Execute(c.Writer, v)
	if err != nil {
		log.Fatal("[FATAL] Execute template: ", err)
		return
	}

	c.Status(200)
}
