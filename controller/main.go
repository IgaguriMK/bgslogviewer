package controller

import (
	"bytes"
	"io"
	"log"

	"github.com/IgaguriMK/bgslogviewer/config"
	"github.com/IgaguriMK/bgslogviewer/model"
	"github.com/IgaguriMK/bgslogviewer/view"
	"github.com/gin-gonic/gin"
)

func MainPage(c *gin.Context) {
	c.Header("Cache-Control", "max-age=600, s-maxage=604800")

	ogp := model.OGP{
		Title:       "BGS Log Viewer",
		Type:        "website",
		Url:         config.BaseUrl(),
		Description: "A viewer for BGS log.",
		HasImage:    true,
		ImageUrl:    config.BaseUrl() + "static/img/ogp/icon_83b3f2.png",
	}

	body := new(bytes.Buffer)

	err := view.Main(body, ogp)
	if err != nil {
		c.String(500, "Internal error")
		log.Fatal("alert: Execute template: ", err)
		return
	}

	c.Status(200)

	io.Copy(c.Writer, body)
}
