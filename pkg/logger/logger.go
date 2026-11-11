package logger

import (
	"fmt"
	"strings"
	"time"
)

type Logger interface {
	Info(msg string)
	Warn(msg string, err error)
	Error(msg string, err error)
}

type Log struct {
}

func New() *Log {
	return &Log{}
}

func (l *Log) Info(msg string) {
	fmt.Println(formatTime(time.Now()), "[INFO]", " "+msg)
}

func (l *Log) Warn(msg string, err error) {
	fmt.Println(formatTime(time.Now()), "[WARN]", " "+msg, err.Error())
}

func (l *Log) Error(msg string, err error) {
	fmt.Println(formatTime(time.Now()), "[ERROR]", " "+msg, err.Error())
}

func formatTime(now time.Time) string {
	time := strings.ReplaceAll(now.UTC().Format(time.RFC3339), "T", " ")
	return strings.ReplaceAll(time, "Z", "")
}
