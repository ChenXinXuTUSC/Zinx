package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type AddOnIntp struct {
	Name string
	Intp func([]interface{}) []interface{}
}

type Logger struct {
	PreAddOnIntps []AddOnIntp
	PstAddOnIntps []AddOnIntp
	LogFile       *os.File
	DumpToFile    bool

	mx sync.Mutex // don't print at the same time
}

func (lp *Logger) Log(format string, args ...interface{}) {
	lp.mx.Lock()
	defer lp.mx.Unlock()
	lp.SlowLog(format, args...)
}

func (lp *Logger) SlowLog(format string, args ...interface{}) {
	var msgs []interface{} = make([]interface{}, 0)

	for _, addon := range lp.PreAddOnIntps {
		msgs = addon.Intp(msgs)
	}
	
	msgs = append(msgs, fmt.Sprintf(format, args...))
	
	for _, addon := range lp.PstAddOnIntps {
		msgs = addon.Intp(msgs)
	}

	fmt.Fprintln(os.Stderr, msgs...)
	if lp.DumpToFile && lp.LogFile != nil {
		fmt.Fprintln(lp.LogFile, msgs...)
	}
}

func (lp *Logger) Dbug(format string, args ...interface{}) {
	lp.Log(LevelMark(lvlDBUG)+" "+format, args...)
}
func (lp *Logger) Info(format string, args ...interface{}) {
	lp.Log(LevelMark(lvlINFO)+" "+format, args...)
}
func (lp *Logger) Warn(format string, args ...interface{}) {
	lp.Log(LevelMark(lvlWARN)+" "+format, args...)
}
func (lp *Logger) Erro(format string, args ...interface{}) {
	lp.Log(LevelMark(lvlERRO)+" "+format, args...)
}
func (lp *Logger) Fatl(format string, args ...interface{}) {
	lp.Log(LevelMark(lvlFATL)+" "+format, args...)
}

func (logger *Logger) RegisterPreAddon(name string, intp func([]interface{}) []interface{}) {
	logger.PreAddOnIntps = append(logger.PreAddOnIntps, AddOnIntp{name, intp})
}
func (logger *Logger) RegisterPstAddon(name string, intp func([]interface{}) []interface{}) {
	logger.PstAddOnIntps = append(logger.PstAddOnIntps, AddOnIntp{name, intp})
}

func NewLogger(logFilePath string, dumpLog bool) Logger {
	file, err := os.OpenFile(
		logFilePath,
		os.O_WRONLY|os.O_CREATE,
		0644,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "try create log file")
		err = os.MkdirAll(filepath.Dir(logFilePath), 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create path dir")
		}
		file, err = os.OpenFile(
			logFilePath,
			os.O_WRONLY|os.O_CREATE,
			0644,
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
	return Logger{
		LogFile: file,
		DumpToFile: dumpLog,
	}
}

type LogLevel int

const (
	lvlFATL LogLevel = iota
	lvlERRO
	lvlWARN
	lvlINFO
	lvlDBUG
)

func LevelMark(lvl LogLevel) string {
	var mark string = "[NONE]"
	switch lvl {
	case lvlDBUG:
		mark = "[DBUG]"
	case lvlINFO:
		mark = "[INFO]"
	case lvlWARN:
		mark = "[WARN]"
	case lvlERRO:
		mark = "[ERRO]"
	case lvlFATL:
		mark = "[FATL]"
	default:
	}
	return mark
}

func PrintDateTime(msgs []interface{}) []interface{} {
	var newMsg string = time.Now().Format("2006-01-02 15:04:05")
	msgs = append(msgs, newMsg)
	return msgs
}

func PrintFileLine(msgs []interface{}) []interface{} {
	_, file, line, ok := runtime.Caller(5)
	var newMsg string = ""
	if !ok {
		newMsg = "unknown stack"
	} else {
		newMsg = fmt.Sprintf("\n    %s:%d", file, line)
	}
	msgs = append(msgs, newMsg)
	return msgs
}

var defaultLogger Logger

func Log(format string, args ...interface{}) {
	defaultLogger.Log(format, args...)
}
func Dbug(format string, args ...interface{}) {
	defaultLogger.Dbug(format, args...)
}
func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}
func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}
func Erro(format string, args ...interface{}) {
	defaultLogger.Erro(format, args...)
}
func Fatl(format string, args ...interface{}) {
	defaultLogger.Fatl(format, args...)
}

func init() {
	defaultLogger = NewLogger(
		filepath.Join(
			"log",
			fmt.Sprintf("%s.log", time.Now().Format("2006-01-02+15:04:05")),
		),
		true,
	)
	defaultLogger.RegisterPreAddon("PrintDateTime", PrintDateTime)
	defaultLogger.RegisterPstAddon("PrintFileLine", PrintFileLine)
}
