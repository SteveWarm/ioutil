package smartnet

import (
    "log"
    "os"
)

type Logger interface {
    Printf(format string, v ...interface{})
}

var DefaultLogger = Logger(log.New(os.Stderr, "smartnet", log.LstdFlags|log.Lshortfile|log.LUTC))

var _ Logger = (*EmptyLogger)(nil)

type EmptyLogger struct {
}

func NewEmptyLogger() *EmptyLogger {
    return &EmptyLogger{}
}

func (this *EmptyLogger) Printf(format string, v ...interface{}) {
}
