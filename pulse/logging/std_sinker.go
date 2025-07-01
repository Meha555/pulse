package logging

import (
	"os"
	"strings"
)

type stdSinker struct {
	formater *LogFormater
	// mtx      sync.Mutex // 不再需要，(os.File).Write底层是加了锁的
}

const (
	// 重置所有颜色设置，恢复默认颜色
	ColorReset = "\033[0m"
	// 前景色为红色
	ColorRed = "\033[31m"
	// 前景色为绿色
	ColorGreen = "\033[32m"
	// 前景色为黄色
	ColorYellow = "\033[33m"
	// 前景色为蓝色
	ColorBlue = "\033[34m"
	// 前景色为灰色
	ColorGray = "\033[90m"
	// 前景色为高亮白色，背景色为红色
	ColorHiRed = "\033[97;41m"
	// 前景色为高亮黄色，背景色为红色
	ColorHiYellowOnRed = "\033[93;41m"
)

func (s *stdSinker) Sink(msg *LogMsg) {
	var builder strings.Builder
	logStr := s.formater.Format(msg)
	var coloredLogStr string
	switch msg.Level {
	case LevelDebug:
		builder.WriteString(ColorGray)
	case LevelWarn:
		builder.WriteString(ColorYellow)
	case LevelError:
		builder.WriteString(ColorRed)
	case LevelFatal:
		builder.WriteString(ColorHiRed)
	case LevelPanic:
		builder.WriteString(ColorHiYellowOnRed)
	}
	builder.WriteString(logStr)
	if msg.Level == LevelPanic {
		builder.WriteByte('\n')
		for _, trace := range msg.Stack {
			builder.WriteString(trace)
			builder.WriteByte('\n')
		}
	}
	builder.WriteString(ColorReset)
	builder.WriteByte('\n')
	coloredLogStr = builder.String()

	// s.mtx.Lock()
	// defer s.mtx.Unlock()
	if msg.Level < LevelError {
		os.Stdout.WriteString(coloredLogStr)
	} else {
		os.Stderr.WriteString(coloredLogStr)
	}
}

func NewStdSinker(format string) *stdSinker {
	return &stdSinker{
		formater: NewLogFormatter(format),
	}
}
