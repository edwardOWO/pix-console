package utils

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

// PixV2Logger defined
type PixV2Logger struct {
	Log *logrus.Logger
}

type Logger interface {
	Debug(where, msg string)
	Info(where, msg string)
	Error(where, msg string)
	Fatal(where, msg string)
}

// InitLogger used to get Logger componnet
func InitLogger(fileName string, level logrus.Level) *PixV2Logger {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	log := logrus.New()
	log.Level = level

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	return &PixV2Logger{
		log}
}

func (logger *PixV2Logger) Debug(where, msg string) {
	logger.Log.WithFields(logrus.Fields{
		"location": where,
	}).Debug(msg)
}

func (logger *PixV2Logger) Info(where, msg string) {
	logger.Log.WithFields(logrus.Fields{
		"location": where,
	}).Info(msg)
}

func (logger *PixV2Logger) Error(where, msg string) {
	logger.Log.WithFields(logrus.Fields{
		"location": where,
	}).Error(msg)
}

func (logger *PixV2Logger) Fatal(where, msg string) {
	logger.Log.WithFields(logrus.Fields{
		"location": where,
	}).Fatal(msg)
}

func Log(logger Logger, level, where, msg string) {
	switch level {
	case "Debug":
		logger.Debug(where, msg)
	case "Info":
		logger.Info(where, msg)
	case "Error":
		logger.Error(where, msg)
	case "Fatal":
		logger.Fatal(where, msg)
	default:
		fmt.Printf("no such log level %s.", level)
	}
}

func Trace() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	_, line := f.FileLine(pc[0])
	return fmt.Sprintf("%v:%v", f.Name(), line)
}
