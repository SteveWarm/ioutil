package watch

import (
    "strconv"
    "time"
)

type Watch time.Time

// 返回当前时间与初始设置的差值
func (me Watch) Sub() time.Duration {
    return time.Now().Sub(time.Time(me))
}

// 返回当前时间与初始设置的差值
// String 是绝大部分日志框架or系统自动反射调用的方法
func (me Watch) String() string {
    return strconv.FormatInt(me.Sub().Milliseconds(), 10)
}


/*
example:
w:=Watch(time.Now())
defer log.Info("us", w)
do something ...
*/