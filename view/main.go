package view

import (
	"bytes"
	"html/template"
	"log"

	"github.com/IgaguriMK/bgslogviewer/model"
)

var mainTemplate *template.Template

func init() {
	var err error

	mainTemplate, err = template.New("main.html.tpl").ParseFiles("template/main.html.tpl")
	if err != nil {
		log.Fatal("[FATAL] Failed parse template: ", err)
	}
}

func Main(res *bytes.Buffer, ogp model.OGP) error {
	err := mainTemplate.Execute(res, ogp)
	if err != nil {
		return err
	}

	return nil
}
