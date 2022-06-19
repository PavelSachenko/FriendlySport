package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

type Logging interface {
	Info(msg string)
	Warning(msg string)
	Error(msg string)
	Debug(msg string)
	Trace(msg string)
	Fatal(msg string)
}

type LogHooks struct {
	LogLevels []logrus.Level
	Writer    io.Writer
}

func (l LogHooks) Levels() []logrus.Level {
	return l.LogLevels
}

func (l LogHooks) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = l.Writer.Write([]byte(line))
	return err
}

type Logger struct {
	*logrus.Logger
}

func (l Logger) Info(msg string) {
	l.Logger.Info(msg)
}

func (l Logger) Warning(msg string) {
	l.Logger.Warning(msg)
}

func (l Logger) Error(msg string) {
	l.Logger.Error(msg)
}

func (l Logger) Debug(msg string) {
	l.Logger.Debug(msg)
}

func (l Logger) Trace(msg string) {
	l.Logger.Trace(msg)
}

func (l Logger) Fatal(msg string) {
	l.Logger.Fatal(msg)
}

func InitLogrusLogger() *Logger {

	l := logrus.New()
	l.SetReportCaller(true)
	l.Info("Init logrus logger")
	l.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
	}
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../..")
	file, err := os.OpenFile(dir+"/logs/workout.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		l.Fatal(err)
	}

	l.SetOutput(ioutil.Discard)
	l.AddHook(&LogHooks{
		Writer: file,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})
	l.AddHook(&LogHooks{
		Writer:    os.Stdout,
		LogLevels: logrus.AllLevels,
	})
	return &Logger{
		Logger: l,
	}
}
