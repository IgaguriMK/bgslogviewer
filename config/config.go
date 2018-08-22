package config

import (
	"os"
	"strconv"
)

var BgsUpdate = 16
var EnableProf = false

func init() {
	if upd, ok := getEnvInt("BLV_BGS_UPDATE"); ok {
		BgsUpdate = upd
	}

	if pf, ok := getEnvBool("BLV_PROFILE"); ok {
		EnableProf = pf
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
