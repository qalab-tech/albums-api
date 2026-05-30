package config

import (
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogger(env string) {
	if env == "development" {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
		Log.SetLevel(logrus.DebugLevel)
	} else {
		Log.SetFormatter(&logrus.JSONFormatter{})
		Log.SetLevel(logrus.InfoLevel)
	}
}
