package controller

import (
	"time"

	"github.com/IgaguriMK/bgslogviewer/config"
	"github.com/gin-gonic/gin"
)

const (
	cacheLocalMin = 120 * time.Second
	cacheLocalMax = 60 * time.Minute
	cacheCdnMin   = 5 * time.Minute
	cacheCdnMax   = 6 * time.Hour
)

func CommonHeader(c *gin.Context) {
	c.Header("X-XSS=Protection", "1; mode=block")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Security-Policy", "default-src 'self'")
}

func cachesLen() (localCache int, cdnCache int) {
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
		return int(cacheLocalMin.Seconds()), int(cacheCdnMin.Seconds())
	}

	// 次のBGS更新まで間近
	if now.Add(time.Hour).After(nextUpdate) {
		return int(cacheLocalMin.Seconds()), int(cacheCdnMin.Seconds())
	}

	// BGS更新付近ではないが、最大キャッシュでは長すぎる
	if now.Add(cacheCdnMax + cacheCdnMin + time.Hour).After(nextUpdate) {
		d := nextUpdate.Sub(now) - time.Hour
		return int(d.Seconds()), int(d.Seconds())
	}

	// 次のBGS更新まで余裕がある
	return int(cacheLocalMax.Seconds()), int(cacheCdnMax.Seconds())
}
