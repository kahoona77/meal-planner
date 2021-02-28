package core

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"os"
)

type IrcService interface {
	Connect(c *Ctx)
	Disconnect(c *Ctx)
}

func InitApp() *Ctx {
	formatter := &logrus.TextFormatter{}
	formatter.ForceColors = true
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(formatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	conf := LoadConfiguration()

	// this connects & tries a simple 'SELECT 1', panics on error
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("sqlite3", conf.DbFile)
	if err != nil {
		panic(err)
	}

	ctx := &Ctx{AppConfig: &conf, Db: db}

	return ctx
}
