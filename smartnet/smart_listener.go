// 多种协议用同一个端口监听
package smartnet

import (
    "bytes"
    "context"
    "io"
    "net"
    "strings"
    "sync"
    "sync/atomic"
    "time"
)

type Config struct {
    PreReadTimeout time.Duration
    ListenConfig   *net.ListenConfig
    Logger         Logger
}

type SmartListener struct {
    network string
    laddr   string
    config  Config
    logger  Logger

    rawListener net.Listener

    socks5Ch       chan net.Conn
    socks5Listener *ChanListener

    httpCh       chan net.Conn
    httpListener *ChanListener

    extListener *ChanListener
    extCh       chan net.Conn

    listenerNums    int64
    connRunableNums int64
    httpFlag        int32
    socks5Flag      int32
    extFlag         int32
    mutex           sync.Mutex
}

func NewSmartListener(network, laddr string, config Config) (*SmartListener, error) {
    if config.ListenConfig == nil {
        config.ListenConfig = &net.ListenConfig{KeepAlive: 5 * time.Minute}
    }

    l, err := config.ListenConfig.Listen(context.Background(), network, laddr)
    if err != nil {
        return nil, err
    }

    s := &SmartListener{
        laddr:       laddr,
        network:     network,
        config:      config,
        rawListener: l,
        extCh:       make(chan net.Conn),
        httpCh:      make(chan net.Conn),
        socks5Ch:    make(chan net.Conn),
    }

    if config.Logger != nil {
        s.logger = config.Logger
    } else {
        s.logger = DefaultLogger
    }

    go s.runable()
    return s, nil
}

func (this *SmartListener) ConnRunables() int64 {
    return atomic.LoadInt64(&this.connRunableNums)
}

func (this *SmartListener) HttpListener() net.Listener {
    return this.initListener(&this.httpListener, this.httpCh, &this.httpFlag)
}

func (this *SmartListener) Socks5Listener() net.Listener {
    return this.initListener(&this.socks5Listener, this.socks5Ch, &this.socks5Flag)
}

func (this *SmartListener) ExtListener() net.Listener {
    return this.initListener(&this.extListener, this.extCh, &this.extFlag)
}

func (this *SmartListener) initListener(l **ChanListener, ch chan net.Conn, v *int32) net.Listener {
    this.mutex.Lock()
    defer this.mutex.Unlock()
    if (*l) == nil {
        *l = NewChanListener(this.rawListener.Addr(), ch)
        atomic.StoreInt32(v, 1)
    }
    return (*l)
}

func (this *SmartListener) HttpEnabled() bool {
    return atomic.LoadInt32(&this.httpFlag) == 1
}

func (this *SmartListener) Socks5Enabled() bool {
    return atomic.LoadInt32(&this.socks5Flag) == 1
}

func (this *SmartListener) ExtEnabled() bool {
    return atomic.LoadInt32(&this.extFlag) == 1
}

func (this *SmartListener) Addr() net.Addr {
    return this.rawListener.Addr()
}

func (this *SmartListener) Close() error {
    this.logger.Printf(this.Addr().String(), "listener close")
    this.rawListener.Close()

    // 强行消费 防止handleConn阻塞
    go func() {
        for {
            select {
            case c, ok := <-this.extCh:
                if !ok {
                    return
                }
                if c != nil {
                    c.Close()
                    this.logger.Printf("force close connection %s <-> %s", c.RemoteAddr(), c.LocalAddr())
                }
            case c, ok := <-this.extCh:
                if !ok {
                    return
                }
                if c != nil {
                    c.Close()
                    this.logger.Printf("force close connection %s <-> %s", c.RemoteAddr(), c.LocalAddr())
                }
            case c, ok := <-this.extCh:
                if !ok {
                    return
                }
                if c != nil {
                    c.Close()
                    this.logger.Printf("force close connection %s <-> %s", c.RemoteAddr(), c.LocalAddr())
                }
            }
        }
    }()

    // 等待所有连接释放
    for {
        if atomic.LoadInt64(&this.listenerNums) > 0 || atomic.LoadInt64(&this.connRunableNums) > 0 {
            time.Sleep(time.Second)
        }
    }

    close(this.extCh)
    close(this.httpCh)
    close(this.socks5Ch)

    if this.HttpEnabled() {
        this.httpListener.Close()
    }

    if this.Socks5Enabled() {
        this.socks5Listener.Close()
    }

    if this.ExtEnabled() {
        this.extListener.Close()
    }

    return nil
}

func (this *SmartListener) runable() {
    atomic.AddInt64(&this.listenerNums, 1)
    defer atomic.AddInt64(&this.listenerNums, -1)

    for {
        c, err := this.rawListener.Accept()
        if err != nil {
            if c != nil {
                panic("BUG: net.Listener returned non-nil conn and non-nil error")
            }
            if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
                this.logger.Printf("Temporary error when accepting new connections: %s", netErr)
                time.Sleep(time.Second)
                continue
            }
            if err != io.EOF && !strings.Contains(err.Error(), "use of closed network connection") {
                this.logger.Printf("Permanent error when accepting new connections: %s", err)
                return
            }
            this.logger.Printf("Some error when accepting new connections: %s", err)
            return
        }
        if c == nil {
            panic("BUG: net.Listener returned (nil, nil)")
        }
        go this.connRunable(c)
    }
}

func (this *SmartListener) connRunable(conn net.Conn) {
    atomic.AddInt64(&this.connRunableNums, 1)
    defer atomic.AddInt64(&this.connRunableNums, -1)

    var buff StackBuffer
    conn.SetReadDeadline(time.Now().Add(this.config.PreReadTimeout))
    // 正常来说这里应该要指定读取字节数，但在当今网络随便一个tcp分片512，第一个包读取不至于太少
    // 所以就读取一次，能读多少是多少 /**n, err := io.ReadAtLeast(conn, buff[:], len(buff))**/
    n, err := conn.Read(buff[:])
    if err != nil {
        // logs.Warn(conn.RemoteAddr(), "closed cause read prefix fail! n:", n, "err:", err)
        err = conn.Close()
        if err != nil {
            // logs.Warn(conn.RemoteAddr(), "close fail! err:", err.Error())
        }
        return
    }

    // 最先判断是不是socks5
    if this.Socks5Enabled() {
        if n >= 3 && isSocks5Request(buff, n) {
            this.socks5Ch <- NewConn(conn, buff, n)
            return
        }
    }

    if this.HttpEnabled() && n >= 8 && isHTTPRequest(buff) {
        this.httpCh <- NewConn(conn, buff, n)
        return
    }

    if this.ExtEnabled() {
        this.extCh <- NewConn(conn, buff, n)
        return
    }

    this.logger.Printf("conn close cause channel not fount %s <-> %s", conn.RemoteAddr(), conn.LocalAddr())
    conn.Close()
}

var httpMethod = [][]byte{
    []byte("GET "),
    []byte("CONNECT "),
    []byte("PUT "),
    []byte("HEAD "),
    []byte("POST "),
    []byte("DELETE "),
    []byte("PATCH "),
    []byte("OPTIONS ")}

func isHTTPRequest(buff StackBuffer) bool {
    for cnt := 0; cnt < len(httpMethod); cnt++ {
        if bytes.HasPrefix(buff[:], []byte(httpMethod[cnt])) {
            return true
        }
    }
    return false
}

func isSocks5Request(buff StackBuffer, n int) bool {
    if buff[0] == 0x05 && buff[1] > 0 && buff[1] < 100 {
        // 检查METHOD是否在合理范围
        for i := 0; i < int(buff[1]) && i < len(buff) && i < (n-2); i++ {
            if !IsInByte(buff[2+i], 0x00, 0x01, 0x02, 0x03, 0x80, 0xFF) {
                return false
            }
        }
        return true
    }
    return false
}

func IsInByte(a byte, b ...byte) bool {
    n := len(b)
    for i := 0; i < n; i++ {
        if a == b[i] {
            return true
        }
    }
    return false
}
