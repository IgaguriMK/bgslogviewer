package controller

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/IgaguriMK/bgslogviewer/apiCaller"
	"github.com/IgaguriMK/bgslogviewer/prof"
	"github.com/IgaguriMK/bgslogviewer/view"
	"github.com/gin-gonic/gin"
)

const (
	systemNameChars = " '+-.0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	systemNameLimit = 40
)

var alloedChars map[byte]bool

func init() {
	alloedChars = make(map[byte]bool)

	for _, c := range []byte(systemNameChars) {
		alloedChars[c] = true
	}
}

func StatPage(c *gin.Context) {
	pf := prof.NewProfiler()
	pf.Start("validation")

	systemName, ok := c.GetQuery("q")
	if !ok {
		c.String(400, "Invalid Url")
		return
	}

	if systemName == "" {
		c.Redirect(301, "/")
		return
	}

	pf.AddParam("q", systemName)

	if !checkSystemName(systemName) {
		c.String(404, "Invalid request")
		pf.End(404)
		log.Printf("info: detect invalid system name %q", systemName)
		return
	}

	pf.Mark("fetchAPI")
	v, stat, err := apiCaller.FetchFactions(systemName)
	if err != nil {
		c.String(500, "Internal error")
		pf.End(500)
		log.Println("error: fetching data error: ", err)
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
		log.Fatal("alert: Execute template: ", err)
		return
	}

	c.Status(200)

	ll, cl := cachesLen()
	c.Header("Cache-Control", fmt.Sprintf("max-age=%d, s-maxage=%d", ll, cl))

	pf.Mark("send")
	io.Copy(c.Writer, body)

	pf.End(200)
}

func checkSystemName(name string) bool {
	if len(name) > systemNameLimit {
		return false
	}

	for _, c := range []byte(name) {
		if !alloedChars[c] {
			return false
		}
	}

	return true
}
