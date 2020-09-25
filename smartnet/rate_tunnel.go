package smartnet

import (
    "fmt"
    "net"
    "os"
    "sync/atomic"
    "time"

    "golang.org/x/time/rate"
)

type StopCallBackFunc func(e1, e2 error)

type CopyConfig struct {
    Limiter      *rate.Limiter
    ReadTimeout  time.Duration // read from src timeout
    WriteTimeout time.Duration // write to dst timeout
}

type TunnelConfig struct {
    A2BConfig    CopyConfig
    B2AConfig    CopyConfig
    StopCallback StopCallBackFunc
}

// 不限速 读写60秒超时
func DefaultTunnelConfig() TunnelConfig {
    return TunnelConfig{
        A2BConfig: CopyConfig{
            Limiter:      nil,
            ReadTimeout:  10 * time.Second,
            WriteTimeout: 10 * time.Second,
        },
        B2AConfig: CopyConfig{
            Limiter:      nil,
            ReadTimeout:  10 * time.Second,
            WriteTimeout: 10 * time.Second,
        },
        StopCallback: nil,
    }
}

// 每秒限速n字节
func RateTunnelConfig(n int) TunnelConfig {
    c := DefaultTunnelConfig()
    c.A2BConfig.Limiter = rate.NewLimiter(rate.Limit(n), n)
    c.B2AConfig.Limiter = rate.NewLimiter(rate.Limit(n), n)
    return c
}

// 带限速功能的通道
type RateTunnel struct {
    a            net.Conn
    b            net.Conn
    config       TunnelConfig
    stopChan     chan int
    stopCallback StopCallBackFunc
    deadNums     int32 // 大于等于2时释放
}

// a, b-连接；a2bConfig-从a读往b写控制配置, b2aConfig-从b读往a写控制配置
func NewRateTunnel(a, b net.Conn, config TunnelConfig) *RateTunnel {
    t := &RateTunnel{a: a, b: b, config: config, stopChan: make(chan int)}
    t.run()
    return t
}

func DefaultRateTunnel(a, b net.Conn) *RateTunnel {
    return NewRateTunnel(a, b, DefaultTunnelConfig())
}

func (this *RateTunnel) Wait() {
    <-this.stopChan
}

// 这里相当于开了3个协程浪费了， 优化到只有2个
func (this *RateTunnel) run() {
    ech := make(chan error, 2)
    go this.copy(this.b, this.a, this.config.B2AConfig, ech)
    go this.copy(this.a, this.b, this.config.A2BConfig, ech)
}

// 速率控制
// limit 控制src到dst的速率 ech-错误输出
func (this *RateTunnel) copy(src, dst net.Conn, config CopyConfig, ech chan error) {
    defer func() {
        if e := recover(); e != nil {
            fmt.Fprintf(os.Stderr,
                "[BUG]func (this *RateTunnel) copy(src, dst net.Conn, config CopyConfig, ech chan<- error) %v", e)
        }
    }()

    defer func() {
        dn := atomic.AddInt32(&this.deadNums, 1)
        if dn <= 1 {
            this.a.Close()
            this.b.Close()
        } else {
            e1 := <-ech
            e2 := <-ech
            close(this.stopChan)
            if this.config.StopCallback != nil {
                this.config.StopCallback(e1, e2)
            }
        }
    }()

    var buff [2000]byte
    var rn, sn int
    var err error
    var re *rate.Reservation

    for {
        if config.ReadTimeout > 0 {
            err = src.SetReadDeadline(time.Now().Add(config.ReadTimeout))
        } else {
            err = src.SetReadDeadline(time.Time{})
        }

        if err != nil {
            ech <- builderror("src.SetReadDeadline", dst, src, err)
            return
        }

        rn, err = src.Read(buff[:])
        if err != nil {
            ech <- builderror("src.Read", dst, src, err)
            return
        }

        if config.Limiter != nil {
            re = config.Limiter.ReserveN(time.Now(), rn)
            if !re.OK() {
                time.Sleep(re.Delay())
            }
        }

        if config.WriteTimeout > 0 {
            err = dst.SetWriteDeadline(time.Now().Add(config.WriteTimeout))
            if err != nil {
                ech <- builderror("dst.SetWriteDeadline", dst, src, err)
                return
            }
        }
        sn, err = dst.Write(buff[0:rn])
        if err != nil {
            ech <- builderror("dst.Write", dst, src, err)
            return
        }

        if sn != rn {
            ech <- builderror("dst.Write incomplete", dst, src, err)
            return
        }
    }
}

func builderror(msg string, src, dst net.Conn, err error) error {
    return fmt.Errorf("%s(%s) -> %s(%s) %s %s",
        src.RemoteAddr().String(), src.LocalAddr().String(),
        dst.RemoteAddr().String(), dst.LocalAddr().String(),
        msg, err.Error(),
    )
}
