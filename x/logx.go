package x

import (
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)
	executionId := xid.New()
	logger.WithFields(logrus.Fields{"executionId": executionId})
	logger.SetReportCaller(true)
	return logger
}