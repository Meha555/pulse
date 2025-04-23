package log

type ILogSinker interface {
	Sink(msg *LogMsg)
}
