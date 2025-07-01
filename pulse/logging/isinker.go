package logging

type ILogSinker interface {
	Sink(msg *LogMsg)
}
