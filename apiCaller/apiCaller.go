package apiCaller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-redis/redis"

	"github.com/IgaguriMK/bgslogviewer/api"
	"github.com/IgaguriMK/bgslogviewer/config"
	"github.com/IgaguriMK/bgslogviewer/model"
)

type ApiStatus int

const (
	Success ApiStatus = iota
	Invalid
	Error
)

const (
	cacheMin     = 5 * time.Minute
	cacheMax     = 3 * time.Hour
	cacheInvalid = 15 * time.Minute
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		DB:   0,
	})
}

var zeroFactions = model.Factions{}
var zeroSystemFactions = api.SystemFactions{}

func FetchFactions(systemName string) (model.Factions, ApiStatus, error) {
	cacheKey := "system\t" + strings.ToLower(systemName)

	bs, err := redisClient.Get(cacheKey).Bytes()
	if err == nil {
		if len(bs) == 0 {
			return zeroFactions, Invalid, nil
		}

		var m model.Factions
		if err := json.Unmarshal(bs, &m); err == nil {
			return m, Success, nil
		}
	}

	res, stat, err := fetchFromEDSM(systemName)
	if err != nil {
		return zeroFactions, Error, err
	}

	if stat == Invalid {
		redisClient.Set(cacheKey, []byte{}, cacheInvalid)
	}

	if stat != Success {
		return zeroFactions, stat, nil
	}

	m := model.FromApiResult(res)
	m.FetchedAt = time.Now()

	bs, err = json.Marshal(m)
	if err != nil {
		return zeroFactions, Error, err
	}
	redisClient.Set(cacheKey, bs, cachesLen())

	return m, Success, nil
}

func cachesLen() time.Duration {
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
		return cacheMin
	}

	// 次のBGS更新まで間近
	if now.Add(time.Hour).After(nextUpdate) {
		return cacheMin
	}

	// BGS更新付近ではないが、最大キャッシュでは長すぎる
	if now.Add(cacheMax + cacheMin + time.Hour).After(nextUpdate) {
		d := nextUpdate.Sub(now) - time.Hour
		return d
	}

	// 次のBGS更新まで余裕がある
	return cacheMax
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
