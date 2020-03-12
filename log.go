package log

import (
	"container/list"
	"fmt"
	"os"
	"time"

	"github.com/op/go-logging"
)

const (
	LoggerName = "log"
)

var (
	Log         = logging.MustGetLogger(LoggerName)
	fileList    = list.New()
	getLogPath  func() string
	getLogLevel func() string
)

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
//     %{id}        Sequence number for log message (uint64).
//     %{pid}       Process id (int)
//     %{time}      Time when log occurred (time.Time)
//     %{level}     Log level (Level)
//     %{module}    Module (string)
//     %{program}   Basename of os.Args[0] (string)
//     %{message}   Message (string)
//     %{longfile}  Full file name and line number: /a/b/c/d.go:23
//     %{shortfile} Final file name element and line number: d.go:23
//     %{color}     ANSI color based on log level
//     %{longpkg}   Full package path, eg. github.com/go-logging
//     %{shortpkg}  Base package path, eg. go-logging
//     %{longfunc}  Full function name, eg. littleEndian.PutUint32
//     %{shortfunc} Base function name, eg. PutUint32
var stdFormat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfile} >%{level:.5s}%{color:reset} - %{message}",
)

var fileFormat = logging.MustStringFormatter(
	"%{time:15:04:05.000} %{shortfile} >%{level:.5s} - %{message}",
)

var auditFormat = logging.MustStringFormatter(
	"%{time:15:04:05.000} - %{message}",
)

// 关闭旧log打开的文件
// newFile本次是否打开了新文件
func closeOldLogFd(newFile bool) {
	expectedFdNum := 0
	if newFile {
		expectedFdNum++
	}
	if fileList.Len() > expectedFdNum {
		element := fileList.Front()
		if element == nil {
			return
		}
		if fp, ok := element.Value.(*os.File); ok {
			fileList.Remove(element)
			time.Sleep(time.Second * 5)
			Log.Notice("start close old log file")
			if err := fp.Close(); err != nil {
				Log.Error("file close err: %v", err)
			}
		} else {
			Log.Error("fd type error")
		}
	}
}

// 如果path不为空，则使用文件记录日志
// 否则使用stdout输出日志
// SetBackend可重复调用
func initLogger(path string, level logging.Level) error {
	if len(path) > 0 {
		fp, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		fileList.PushBack(fp)
		fileBackend := logging.NewLogBackend(fp, "", 1)
		fileFormatter := logging.NewBackendFormatter(fileBackend, fileFormat)

		fileB := logging.AddModuleLevel(fileFormatter)
		fileB.SetLevel(level, "")
		logging.SetBackend(fileB)
	} else {
		stdBackend := logging.NewLogBackend(os.Stdout, "", 1)
		stdFormatter := logging.NewBackendFormatter(stdBackend, stdFormat)
		stdB := logging.AddModuleLevel(stdFormatter)
		stdB.SetLevel(level, "")
		logging.SetBackend(stdB)
	}
	go closeOldLogFd(len(path) > 0)
	return nil
}

func InitLogger(getLogPathFunc func() string, getLogLevelFunc func() string) {
	getLogPath = getLogPathFunc
	getLogLevel = getLogLevelFunc
	reloadLog()
}

func reloadLog() {
	logPath := getLogPath()
	logLevel := getLogLevel()
	level, err := logging.LogLevel(logLevel)
	if err != nil {
		level = logging.INFO
	}
	if err := initLogger(logPath, level); err != nil {
		fmt.Printf("init logger %s error: %v\n", logPath, err)
		os.Exit(1)
	}
	Log = logging.MustGetLogger(LoggerName)
}
