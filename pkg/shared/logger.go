package shared

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var detailedLoggingLock sync.Once
var logConfig LogConfig

func InitLogger() {
	detailedLoggingLock.Do(func() {
		logConfig = configInstance.LogConfig
	})

	if logConfig.UseJsonFormatter {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		textFormatter := new(logrus.TextFormatter)
		textFormatter.ForceColors = true
		textFormatter.FullTimestamp = true
		textFormatter.TimestampFormat = "2006-01-02 12:13:15.000"
		textFormatter.ForceQuote = true
		textFormatter.PadLevelText = true
		textFormatter.DisableLevelTruncation = true
		logrus.SetFormatter(textFormatter)
	}

	if logConfig.DetailedLogs {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(false)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetOutput(os.Stdout)
}

func HandleError(err error, shouldPanic bool) {
	HandleErrorWithMsg(err, shouldPanic, "")
}

func HandleErrorWithMsg(err error, shouldPanic bool, msg string) {
	if err != nil {
		if msg == "" {
			msg = "An error occurred"
		}

		logrus.WithError(err).Error(msg)

		if shouldPanic {
			panic(err)
		}
	}
}
