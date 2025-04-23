package log

import (
	"fmt"
	"os"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal // 记录日志并额外执行 os.Exit(1)
	LevelPanic // 记录日志并额外执行一条 panic() 语句
)

var LevelStrs = map[int]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
	LevelPanic: "PANIC",
}

// 预先定义一个全局默认的Logger实例，方便使用
var Log *Logger

func init() {
	Log = NewStdLogger(LevelInfo, "DEFAULT", "[%t] [%c %l] [%f:%C:%L:%g] %m", false)
}

// NOTE 如果LoggerBase是小写的，则子包也无法看到这个符号

type Logger struct {
	baseLevel int
	category  string
	sinker    ILogSinker
	logCh     chan *LogMsg
}

func (l *Logger) enable(level int) bool {
	return level >= l.baseLevel
}

func NewStdLogger(level int, category, format string, async bool) *Logger {
	l := &Logger{
		baseLevel: level,
		category:  category,
		sinker:    NewStdSinker(format),
		logCh:     nil,
	}

	if async {
		l.logCh = make(chan *LogMsg, 10)
		go func() {
			for msg := range l.logCh {
				l.sinker.Sink(msg)
			}
		}()
	}

	return l
}

func NewFileLogger(level int, category, format, filePath, fileName string, maxLogSize int64, async bool) *Logger {
	l := &Logger{
		baseLevel: level,
		category:  category,
		sinker:    NewFileSinker(format, filePath, fileName, maxLogSize),
		logCh:     nil,
	}

	if async {
		l.logCh = make(chan *LogMsg, 10)
		go func() {
			for msg := range l.logCh {
				l.sinker.Sink(msg)
			}
		}()
	}

	return l
}

func (l *Logger) doLog(level int, content string, callDepth int, traceSkip int) {
	if l.enable(level) {
		msg := NewMsg(level, l.category, content).WithCallDepth(callDepth)
		if traceSkip >= 0 {
			msg.WithStack(traceSkip)
		}
		if l.logCh == nil {
			l.sinker.Sink(msg)
		} else {
			msg.WithFile(3)
			l.logCh <- msg
		}
	}
}

func (l *Logger) Log(level int, content string) {
	l.doLog(level, content, 6, -1)
}

func (l *Logger) Debug(content string) {
	l.doLog(LevelDebug, content, 6, -1)
}

func (l *Logger) Info(content string) {
	l.doLog(LevelInfo, content, 6, -1)
}

func (l *Logger) Warn(content string) {
	l.doLog(LevelWarn, content, 6, -1)
}

func (l *Logger) Error(content string) {
	l.doLog(LevelError, content, 6, -1)
}

func (l *Logger) Fatal(content string) {
	l.doLog(LevelFatal, content, 6, -1)
	os.Exit(1)
}

// Panic 为了排除log包的调用栈，不得不自行构造LogMsg，然后传入sinker.Sink()
func (l *Logger) Panic(content string) {
	l.doLog(LevelPanic, content, 6, 3)
	panic(content)
}

func (l *Logger) Logf(level int, format string, v ...interface{}) {
	l.doLog(level, fmt.Sprintf(format, v...), 6, -1)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.doLog(LevelDebug, fmt.Sprintf(format, v...), 6, -1)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.doLog(LevelInfo, fmt.Sprintf(format, v...), 6, -1)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.doLog(LevelWarn, fmt.Sprintf(format, v...), 6, -1)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.doLog(LevelError, fmt.Sprintf(format, v...), 6, -1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.doLog(LevelFatal, fmt.Sprintf(format, v...), 6, -1)
	os.Exit(1)
}

// Panicf 为了排除log包的调用栈，不得不自行构造LogMsg，然后传入sinker.Sink()
func (l *Logger) Panicf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v...)
	l.doLog(LevelPanic, content, 6, 3)
	panic(content)
}

func (l *Logger) Level() int {
	return l.baseLevel
}

func (l *Logger) SetLevel(level int) {
	// 不对baseLevel的修改进行加锁，这里对一致性要求不高
	l.baseLevel = level
}
