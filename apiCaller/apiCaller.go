package apiCaller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/IgaguriMK/bgslogviewer/api"
	"github.com/IgaguriMK/bgslogviewer/model"
)

type ApiStatus int

const (
	Success ApiStatus = iota
	Invalid
	Error
)

var zeroFactions = model.Factions{}
var zeroSystemFactions = api.SystemFactions{}

func FetchFactions(systemName string) (model.Factions, ApiStatus, error) {
	res, stat, err := fetchFromEDSM(systemName)
	if err != nil {
		return zeroFactions, Error, err
	}

	if stat != Success {
		return zeroFactions, stat, nil
	}

	return model.FromApiResult(res), Success, nil
}

func fetchFromEDSM(systemName string) (api.SystemFactions, ApiStatus, error) {
	params := url.Values{}
	params.Add("systemName", systemName)
	params.Add("showHistory", "1")

	url := "https://www.edsm.net/api-system-v1/factions?" + params.Encode()
	resp, err := http.Get(url)
	if err != nil {
		return zeroSystemFactions, Error, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return zeroSystemFactions, Error, err
	}

	if bytes.Equal(buf.Bytes(), []byte("{}")) {
		return zeroSystemFactions, Invalid, nil
	}

	var res api.SystemFactions
	err = json.NewDecoder(buf).Decode(&res)
	if err != nil {
		return zeroSystemFactions, Error, err
	}

	return res, Success, nil
}
