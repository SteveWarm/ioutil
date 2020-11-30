package syscmd

import (
    "fmt"
    "sync/atomic"
)

var _no_run int32 = 0 // 不运行

func NoRunSwitch(on bool) {
    if on {
        atomic.StoreInt32(&_no_run, 1)
    } else {
        atomic.StoreInt32(&_no_run, 0)
    }
}

func NoRunEnabled() bool {
    return atomic.LoadInt32(&_no_run) == 1
}

var _log_on int32 = 0

func LogSwitch(on bool) {
    if on {
        atomic.StoreInt32(&_log_on, 1)
    } else {
        atomic.StoreInt32(&_log_on, 0)
    }
}

func LogEnabled() bool {
    return atomic.LoadInt32(&_log_on) == 1
}

func log(args ...interface{}) {
    if LogEnabled() {
        fmt.Println(args...)
    }
}

func logf(format string, args ...interface{}) {
    if LogEnabled() {
        fmt.Printf(format, args...)
    }
}
