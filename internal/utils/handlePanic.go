package utils

import (
	"os"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

// RecoverAndLogPanic recovers from panic, logs it, and exits
func RecoverAndLogPanic(logger *logrus.Logger) {
	if r := recover(); r != nil {
		logger.Errorf("Panic occurred: %v", r)
		logger.Debugf("Stacktrace:\n%s", debug.Stack())
		os.Exit(1)
	}
}
