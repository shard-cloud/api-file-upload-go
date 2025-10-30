package logger

import (
	"os"
	"github.com/sirupsen/logrus"
)

func New(level string) *logrus.Logger {
	log := logrus.New()
	
	// Set log level
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
	
	// Set output
	log.SetOutput(os.Stdout)
	
	// Set formatter
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	
	return log
}
