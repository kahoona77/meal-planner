package core

import (
	"os"
	"strconv"
)

type AppConfig struct {
	IsDev    bool
	Port     string
	BasePath string
	DbFile   string
}

func LoadConfiguration() AppConfig {
	conf := AppConfig{}

	conf.IsDev, _ = strconv.ParseBool(os.Getenv("DEV_MODE"))

	conf.Port = os.Getenv("PORT")
	if conf.Port == "" {
		conf.Port = "8080"
	}

	conf.BasePath = os.Getenv("BASE_PATH")
	if conf.BasePath == "" {
		conf.BasePath = "/meal-planner"
	}

	conf.DbFile = os.Getenv("DB_FILE")
	if conf.DbFile == "" {
		conf.DbFile = "meal-planner.sqlite"
	}

	return conf
}
