package core

import (
	"log/slog"
	"os"
)

type AppConfig struct {
	Port      string
	BasePath  string
	DbFile    string
	LogLevel  slog.Level
	StaticDir string
}

func LoadConfiguration() AppConfig {
	conf := AppConfig{}

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

	// STATIC_DIR serves the static files from disk instead of the copy in the binary (development)
	conf.StaticDir = os.Getenv("STATIC_DIR")

	// LOG_LEVEL is one of debug, info, warn, error; defaults to info
	if err := conf.LogLevel.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err != nil {
		conf.LogLevel = slog.LevelInfo
	}

	return conf
}
