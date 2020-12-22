package weixinmp

type EmptyLogger struct {
}

func (me EmptyLogger) Error(...interface{}) {}
func (me EmptyLogger) Info(...interface{})  {}
func (me EmptyLogger) Warn(...interface{})  {}
func (me EmptyLogger) Debug(...interface{}) {}

type Logger interface {
    Error(...interface{})
    Info(...interface{})
    Warn(...interface{})
    Debug(...interface{})
}

var logger Logger = EmptyLogger{}

// 非线程安全：要保证在程序最开始调用
func SetLogger(l Logger) {
    logger = l
}
