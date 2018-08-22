package config

import (
	"os"
	"strconv"
)

var BgsUpdate = 16
var EnableProf = false
var RedisHost = "redis"
var RedisPort = 6379

func init() {
	if upd, ok := getEnvInt("BLV_BGS_UPDATE"); ok {
		BgsUpdate = upd
	}

	if pf, ok := getEnvBool("BLV_PROFILE"); ok {
		EnableProf = pf
	}

	if h, ok := os.LookupEnv("REDIS_HOST"); ok {
		RedisHost = h
	}

	if p, ok := getEnvInt("REDIS_PORT"); ok {
		RedisPort = p
	}
}

func getEnvInt(key string) (int, bool) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return 0, false
	}

	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false
	}

	return int(n), true
}

func getEnvBool(key string) (bool, bool) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return false, false
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, false
	}

	return b, true
}
