package logging

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/siddontang/go-log/loggers"
)

type binlogLogger struct{}

var (
	_ loggers.Advanced = (*binlogLogger)(nil)
)

func NewBinlogLogger() loggers.Advanced {
	return &binlogLogger{}
}

func (*binlogLogger) Fatal(args ...interface{}) {
	slog.Error(fmt.Sprint(args...))
	os.Exit(1)
}

func (*binlogLogger) Fatalf(format string, args ...interface{}) {
	slog.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (*binlogLogger) Fatalln(args ...interface{}) {
	slog.Error(fmt.Sprintln(args...))
	os.Exit(1)
}

func (*binlogLogger) Panic(args ...interface{}) {
	slog.Error(fmt.Sprintln(args...))
	panic(args)
}

func (*binlogLogger) Panicf(format string, args ...interface{}) {
	slog.Error(fmt.Sprintf(format, args...))
	panic(fmt.Errorf(format, args...))
}

func (*binlogLogger) Panicln(args ...interface{}) {
	slog.Error(fmt.Sprintln(args...))
	panic(args)
}

func (*binlogLogger) Print(args ...interface{}) {
	slog.Info(fmt.Sprint(args...))
}

func (*binlogLogger) Printf(format string, args ...interface{}) {
	slog.Info(fmt.Sprintf(format, args...))
}

func (*binlogLogger) Println(args ...interface{}) {
	slog.Info(fmt.Sprintln(args...))
}

func (*binlogLogger) Debug(args ...interface{}) {
	slog.Debug(fmt.Sprint(args...))
}

func (*binlogLogger) Debugf(format string, args ...interface{}) {
	slog.Debug(fmt.Sprintf(format, args...))
}

func (*binlogLogger) Debugln(args ...interface{}) {
	slog.Debug(fmt.Sprintln(args...))
}

func (*binlogLogger) Error(args ...interface{}) {
	slog.Error(fmt.Sprint(args...))
}

func (*binlogLogger) Errorf(format string, args ...interface{}) {
	slog.Error(fmt.Sprintf(format, args...))
}

func (*binlogLogger) Errorln(args ...interface{}) {
	slog.Error(fmt.Sprintln(args...))
}

func (*binlogLogger) Info(args ...interface{}) {
	slog.Info(fmt.Sprint(args...))
}

func (*binlogLogger) Infof(format string, args ...interface{}) {
	slog.Info(fmt.Sprintf(format, args...))
}

func (*binlogLogger) Infoln(args ...interface{}) {
	slog.Info(fmt.Sprintln(args...))
}

func (*binlogLogger) Warn(args ...interface{}) {
	slog.Warn(fmt.Sprint(args...))
}

func (*binlogLogger) Warnf(format string, args ...interface{}) {
	slog.Warn(fmt.Sprintf(format, args...))
}

func (*binlogLogger) Warnln(args ...interface{}) {
	slog.Warn(fmt.Sprintln(args...))
}
