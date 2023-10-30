package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

type Service struct {
	logLevel          string
	logFilepath       string
	consoleLogEnabled bool
}

func NewService(
	logLevel string,
	logFilepath string,
	consoleLogEnabled bool,
) *Service {
	return &Service{
		logLevel:          logLevel,
		logFilepath:       logFilepath,
		consoleLogEnabled: consoleLogEnabled,
	}
}

func (s *Service) Init() error {
	logrus.New()
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339Nano,
		DisableHTMLEscape: false,
		PrettyPrint:       true,
	})

	logLevel, err := logrus.ParseLevel(s.logLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(logLevel)

	logFile, err := os.OpenFile(s.logFilepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	var multiWriter io.Writer
	if s.consoleLogEnabled {
		multiWriter = io.MultiWriter(logFile, os.Stdout)
	} else {
		multiWriter = io.MultiWriter(logFile)
	}
	logrus.SetOutput(multiWriter)
	return nil
}
