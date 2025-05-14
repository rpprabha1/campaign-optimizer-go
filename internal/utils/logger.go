package utils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func NewLogger(serviceName string) *logrus.Logger {
	// Define the logs directory and file path
	logDir := filepath.Join("..", "..", "logs")
	logFileName := filepath.Join(logDir, serviceName + ".log")

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0777); err != nil {
		panic("Unable to create log directory: " + err.Error())
	}

	// Create or open the log file
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		panic("Unable to open log file: " + err.Error())
	}

	// Initialize logger
	logger := logrus.New()

	// Output to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(multiWriter)

	// Log format
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Minimum level to log
	logger.SetLevel(logrus.InfoLevel)

	return logger
}
