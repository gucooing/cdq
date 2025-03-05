package logger

import (
	"fmt"
	"io"
	"time"
)

type Log struct {
	level LogLevel
	out   io.Writer
}

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func NewLog(level LogLevel, out io.Writer) *Log {
	return &Log{
		level: level,
		out:   out,
	}
}

func (x *Log) Debug(msg string, param ...interface{}) {
	if x.level > LevelDebug {
		return
	}

	x.log(LevelDebug, msg, param...)
}

func (x *Log) Info(msg string, v ...interface{}) {
	if x.level > LevelInfo {
		return
	}

	x.log(LevelInfo, msg, v...)
}

func (x *Log) Warn(msg string, v ...interface{}) {
	if x.level > LevelWarn {
		return
	}

	x.log(LevelWarn, msg, v...)
}

func (x *Log) Error(msg string, v ...interface{}) {
	if x.level > LevelError {
		return
	}

	x.log(LevelError, msg, v...)
}

func (x *Log) log(level LogLevel, msg string, v ...interface{}) {
	colorCode := "0"
	switch level {
	case LevelError:
		colorCode = "31" // 红色
	case LevelWarn:
		colorCode = "33" // 黄色
	case LevelInfo:
		colorCode = "32" // 绿色
	case LevelDebug:
		colorCode = "37" // 灰色
	}

	message := fmt.Sprintf(msg, v...)
	logLine := fmt.Sprintf("\033[%sm[%s][%s] %s\033[0m\n", colorCode, level, time.Now().Format("2006-01-02 15:04:05.000 MST"), message)
	fmt.Print(logLine)
	if x.out != nil {
		x.out.Write([]byte(logLine))
	}
}
