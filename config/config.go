package config

import (
	"fmt"
	"os"
	"strconv"
)

var (
	Debug      = false
	BgsUpdate  = 16
	EnableProf = false
	RedisHost  = "redis"
	RedisPort  = 6379
	Protocol   = "http"
	HostName   = ""
)

func init() {

	if d, ok := getEnvBool("BLV_DEBUG"); ok {
		Debug = d
	}

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

	if h, ok := getEnvString("BLV_HOSTNAME"); ok {
		HostName = h
	} else {
		panic("must set BLV_HOSTNAME")
	}

	if p, ok := getEnvString("BLV_PROTO"); ok {
		Protocol = p
	}
}

func BaseUrl() string {
	return fmt.Sprintf("%s://%s/", Protocol, HostName)
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

func getEnvString(key string) (string, bool) {
	v, ok := os.LookupEnv(key)
	return v, ok
}
