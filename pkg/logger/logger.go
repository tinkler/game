package logger

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"sfs.ink/liang/game/pkg/status"
)

var (
	mu    sync.Mutex
	level LogLevel
)

type LogLevel int

const (
	_ LogLevel = iota
	LErr
	LWarn
	LDebug
	LInfo
)

type ConsoleType int

const (
	_ ConsoleType = iota + 29
	ConsoleGray
	ConsoleRed
	ConsoleGreen
	ConsoleYellow
	ConsoleBlue
)

var DefaultLogger = &defaultLogger

var defaultLogger = Logger{}

type Logger struct {
	s status.Status
}

func NewLogger(s status.Status) *Logger {
	return &Logger{s: s}
}

func format(f interface{}, v ...interface{}) string {
	var msg string
	switch f := f.(type) {
	case string:
		msg = f
		if len(v) == 0 {
			return msg
		}
		if !strings.Contains(msg, "%") {
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	if len(v) > 0 {
		return fmt.Sprintf(msg, v...)
	}
	return fmt.Sprint(msg)
}

func writeRuntimeMsg(dep int) string {
	_, file, line, _ := runtime.Caller(dep)
	return fmt.Sprintf("\x1b[1;%dm%s:%d\x1b[0m", ConsoleGray, file, line)
}

func writeMsg(level LogLevel, msg string) string {
	ts := time.Now().Format("2006-01-02 15:04:05")
	switch level {
	case LInfo:
		msg = fmt.Sprintf("%s \x1b[0;%dm[日志]\x1b[0m %s", ts, ConsoleBlue, msg)
	case LDebug:
		msg = fmt.Sprintf("%s \x1b[0;%dm[调试]\x1b[0m %s", ts, ConsoleGreen, msg)
	case LWarn:
		msg = fmt.Sprintf("%s \x1b[0;%dm[警告]\x1b[0m %s", ts, ConsoleYellow, msg)
	case LErr:
		msg = fmt.Sprintf("%s \x1b[0;%dm[错误]\x1b[0m %s", ts, ConsoleRed, msg)
	}
	return msg
}

func Info(f interface{}, v ...interface{}) {
	os.Stdout.Write(append([]byte(writeMsg(LInfo, format(f, v...))), '\n'))
}

func Warn(f interface{}, v ...interface{}) {
	os.Stdout.Write(
		append(
			append([]byte(writeRuntimeMsg(2)), '\n'),
			append([]byte(writeMsg(LWarn, format(f, v...))), '\n')...))
}

// Info 普通日志打印
func (l *Logger) Info(f interface{}, v ...interface{}) {
	os.Stdout.Write(
		append(
			[]byte(fmt.Sprintf("Type: %s, ID: %s >>>>>\n", l.s.TypeName(), l.s.ID())),
			append([]byte(writeMsg(LInfo, format(f, v...))), []byte("\n<<<<<\n")...)...,
		),
	)
}

func (l *Logger) Warn(f interface{}, v ...interface{}) {
	bys := bytes.NewBuffer([]byte(fmt.Sprintf("Type: %s, ID: %s >>>>>", l.s.TypeName(), l.s.ID())))
	bys.WriteByte('\n')
	bys.Write([]byte(writeMsg(LWarn, format(f, v...))))
	bys.Write([]byte("\n<<<<<\n"))
	os.Stdout.Write(bys.Bytes())
}

func (l *Logger) Error(f interface{}, v ...interface{}) {
	bys := bytes.NewBuffer([]byte(writeRuntimeMsg(2)))
	bys.WriteByte('\n')
	bys.Write([]byte(fmt.Sprintf("Type: %s, ID: %s >>>>>", l.s.TypeName(), l.s.ID())))
	bys.WriteByte('\n')
	bys.Write([]byte(writeMsg(LErr, format(f, v...))))
	bys.Write([]byte("\n<<<<<\n"))
	os.Stdout.Write(bys.Bytes())
}
