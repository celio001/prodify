package config

import (
	"os"
	"strconv"
)

var config = map[string]string{
	//postgress
	"USER_POST":     "USER_POST",
	"HOST_POST":     "HOST_POST",
	"PORT_POST":     "PORT_POST",
	"PASSWORD_POST": "PASSWORD_POST",
	"DB_NAME_POST":  "DB_NAME_POST",
}

func GetString(k string) string {
	v := os.Getenv(k)
	if v == "" {
		return config[k]
	}
	return v
}

func GetInt(k string) int {
	v := GetString(k)
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}

	return i
}
