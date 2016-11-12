package logger

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
	"strings"
	"path/filepath"
)

const (
	LOG_LEVEL_DEBUG = iota // 0
	LOG_LEVEL_INFO
	LOG_LEVEL_NOTICE
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR
)

var LOG_LEVEL_MAP map[int]string = map[int]string{
	LOG_LEVEL_DEBUG : "DEBUG",
	LOG_LEVEL_INFO : "INFO",
	LOG_LEVEL_NOTICE : "NOTICE",
	LOG_LEVEL_WARN : "WARN",
	LOG_LEVEL_ERROR : "ERROR",
}

type Logger struct {
	id string
	appName string
	logDir string
	logLevel int
	fileLock sync.Mutex
	lineLock sync.Mutex
	file *os.File
	commonFields string
}

var _logger *Logger = &Logger{}

func Init(appName string, logDir string, logLevel int){
	_logger.appName = appName
	_logger.logDir = strings.TrimRight(logDir, "/") + "/"
	_logger.logLevel = logLevel

	reopen()
}

func reopen(){
	tm := time.Unix(time.Now().Unix(), 0)
	h := tm.Format("2006010215")

	logfile := fmt.Sprintf("%s/%s.log.%s",
		_logger.logDir,
		_logger.appName,
		h)
	
	//file does not exist
	_, ferr := os.Stat(logfile)
	if ferr != nil {
		_logger.fileLock.Lock()
		defer _logger.fileLock.Unlock()
		if _logger.file != nil {
			_logger.file.Close()
		}
		f, ferr := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if ferr != nil {
			fmt.Println("open log file failed", logfile)
			return
		}
		_logger.file = f
		return
	}
	//file exist,but has not file handler
	if _logger.file == nil {
		_logger.fileLock.Lock()
		defer _logger.fileLock.Unlock()
		f, ferr := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if ferr != nil {
			fmt.Println("open log file failed", logfile)
			return
		}
		_logger.file = f
		return
	}
}

func WithCommonFields(fields map[string]string){
	if len(fields) == 0 {
		return
	}
	
	for k,v := range(fields){
		_logger.commonFields += k+"=["+v+"] "
	}
	_logger.commonFields = strings.TrimRight(_logger.commonFields, " ")
}

func SetLogId(id string){
	_logger.id = id
}

func GetLogId(id string) string{
	return _logger.id
}


func Debug(format string, args ...interface{}){
	output(LOG_LEVEL_DEBUG, format, args...)
}

func Info(format string, args ...interface{}){
	output(LOG_LEVEL_INFO, format, args...)
}

func Warn(format string, args ...interface{}){
	output(LOG_LEVEL_WARN, format, args...)
}

func Error(format string, args ...interface{}){
	output(LOG_LEVEL_ERROR, format, args...)
}

func Notice(format string, args ...interface{}){
	output(LOG_LEVEL_NOTICE, format, args...)
}

func output(level int, format string, args ...interface{}){
	if _logger.logLevel > level {
		return
	}

	tm := time.Unix(time.Now().Unix(), 0)
	s := tm.Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}else{
		file = filepath.Base(file)
	}

	entryfmt := "%s %s %s %s:%d %s %s %s\n"
	entrystr := fmt.Sprintf(entryfmt, _logger.appName, s, LOG_LEVEL_MAP[level],
					file, line, _logger.id, _logger.commonFields, msg)

	reopen()
	
	//_logger.lineLock.Lock()
	//defer _logger.lineLock.Unlock()
	
	_, err := _logger.file.WriteString(entrystr)
	if err != nil {
		fmt.Println("log write file error", err.Error())
	}
}

