package utils

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger(serviceName string) *logrus.Logger {
	// Create log file with timestamp-based naming
	logFileName := serviceName +  ".log"
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Unable to open log file: " + err.Error())
	}

	logger := logrus.New()

	// Output to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(multiWriter)

	// Log format
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Minimum level to log
	logger.SetLevel(logrus.InfoLevel)

	return logger
}