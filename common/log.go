package common

import (
	"regexp"
	"strings"

	"github.com/FoxComm/FoxComm/logger"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	ctx logger.LogContext
}

func NewLogger(funcRegxp, fileRegxp, packRegxp string, skipStack int) (*Logger, error) {
	funcRegexp, err := regexp.Compile(funcRegxp)
	if err != nil {
		return nil, err
	}
	fileRegexp, err := regexp.Compile(fileRegxp)
	if err != nil {
		return nil, err
	}
	packRegexp, err := regexp.Compile(packRegxp)
	if err != nil {
		return nil, err
	}

	return &Logger{ctx: logger.LogContext{
		FuncRegexp: funcRegexp,
		FileRegexp: fileRegexp,
		PackRegexp: packRegexp,
		SkipStack:  skipStack,
	},
	}, nil
}

func (l Logger) SetSeverity(s string) error {
	lvlStr := strings.ToLower(s)
	level, err := logrus.ParseLevel(lvlStr)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	return nil
}

func (l Logger) GetSeverity() string {
	return logrus.GetLevel().String()
}

func (l Logger) Infof(format string, args ...interface{}) {
	logger.LogWithContext(format, logger.DEBUG, l.ctx, args...)
}
func (l Logger) Warningf(format string, args ...interface{}) {
	logger.LogWithContext(format, logger.WARN, l.ctx, args...)
}
func (l Logger) Errorf(format string, args ...interface{}) {
	logger.LogWithContext(format, logger.ERROR, l.ctx, args...)
}
func (l Logger) Fatalf(format string, args ...interface{}) {
	logger.LogWithContext(format, logger.FATAL, l.ctx, args...)
}
