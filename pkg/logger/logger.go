package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"runtime"
)

type writerHook struct {
	Writer   []io.Writer
	LogLevel []logrus.Level
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() *Logger {
	return &Logger{e}
}

func (l *Logger) GetLoggerWithField(key string, value interface{}) *Logger {
	return &Logger{l.WithField(key, value)}
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}

	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevel
}

func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	})

	err := os.Mkdir("logs", 0644)
	if err != nil {
		log.Fatalf("failed to make dir: %v", err)
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:   []io.Writer{allFile, os.Stdout},
		LogLevel: logrus.AllLevels,
	})
}
