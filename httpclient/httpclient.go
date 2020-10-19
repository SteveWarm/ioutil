// 支持http http2 ssl认证 ssl双向认证 Cookies的HttpClient封装
package httpclient

import (
    "context"
    "crypto/tls"
    "golang.org/x/net/http2"
    "net"
    "net/http"
    "net/http/cookiejar"
    "os"
    "syscall"
    "time"
)

const (
    DefaultDialTimeout         = 3 * time.Second
    DefaultTimeout             = 10 * time.Second
    DefaultTLSHandshakeTimeout = 3 * time.Second
    DefaultIdleConnTimeout     = 50 * time.Second
    DefaultMaxConnsPerHost     = 30
    DefaultMaxIdleConns        = 1000
    DefaultMaxIdleConnsPerHost = 1

    ProtoH1 = "http/1.1"
    ProtoH2 = "h2"
)

type Options struct {
    JarEnabled           bool              // 是否启用 cookies
    Timeout              time.Duration     // 默认10秒
    DialTimeout          time.Duration     //连接超时时间 默认3秒
    TLSHandshakeTimeout  time.Duration     // 默认3秒
    IdleConnTimeout      time.Duration     // 默认50秒
    MaxConnsPerHost      int               // 每个域名最多连接数量
    MaxIdleConns         int               // 最多允许的空闲连接数量
    MaxIdleConnsPerHost  int               // 每个域名最多允许的空闲连接数量
    WriteBufferSize      int               // 读写缓冲区 根据业务调整合适 如果都是小报文响应 适当调小
    ReadBufferSize       int               // 读写缓冲区 根据业务调整合适 如果都是小报文请求 适当调小
    InsecureSkipVerify   bool              // 默认false
    CertFilePath         string            // 从文件中加载证书 跟Certificates合并
    KeyFilePath          string            // 从文件中加载证书 跟Certificates合并
    Certificates         []tls.Certificate // 双向认证时使用
    ProxyFromEnvironment bool              // 是否用环境代理 默认不用代理
    // NextProtos           []string          // 协商使用的协议 由http和http2包内不自动设置 ProtoH1  ProtoH2
}

func NewHttpClient(option Options) (*http.Client, error) {
    if option.Timeout == 0 {
        option.Timeout = DefaultTimeout
    }

    if option.DialTimeout == 0 {
        option.DialTimeout = DefaultDialTimeout
    }

    if option.TLSHandshakeTimeout == 0 {
        option.TLSHandshakeTimeout = DefaultTLSHandshakeTimeout
    }

    if option.IdleConnTimeout == 0 {
        option.IdleConnTimeout = DefaultIdleConnTimeout
    }

    if option.MaxIdleConns == 0 {
        option.MaxIdleConns = DefaultMaxIdleConns
    }

    if option.MaxIdleConnsPerHost == 0 {
        option.MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
    }

    if option.MaxConnsPerHost == 0 {
        option.MaxConnsPerHost = DefaultMaxConnsPerHost
    }

    if option.CertFilePath != "" && option.KeyFilePath != "" {
        cert, err := tls.LoadX509KeyPair(option.CertFilePath, option.KeyFilePath)
        if err != nil {
            return nil, err
        }
        option.Certificates = append(option.Certificates, cert)
    }

    tlsCfg := &tls.Config{
        InsecureSkipVerify: option.InsecureSkipVerify,
        Certificates:       option.Certificates,
        //  NextProtos: no need to set. auto set int http and https package
    }

    tr := &http.Transport{
        DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
            return net.DialTimeout(network, addr, option.DialTimeout)
        },
        TLSHandshakeTimeout: option.TLSHandshakeTimeout,

        MaxIdleConns:        option.MaxIdleConns,
        MaxIdleConnsPerHost: option.MaxIdleConnsPerHost,
        MaxConnsPerHost:     option.MaxConnsPerHost,
        IdleConnTimeout:     option.IdleConnTimeout,
        WriteBufferSize:     option.WriteBufferSize,
        ReadBufferSize:      option.ReadBufferSize,
        TLSClientConfig:     tlsCfg,
    }

    if option.ProxyFromEnvironment {
        tr.Proxy = http.ProxyFromEnvironment
    }

    err := http2.ConfigureTransport(tr)
    if err != nil {
        return nil, err
    }

    cli := &http.Client{
        Transport: tr,
        Timeout:   option.Timeout,
    }

    if option.JarEnabled {
        jar, err := cookiejar.New(nil)

        if err != nil {
            return nil, err
        }

        cli.Jar = jar
    }

    // fmt.Println(tlsCfg.NextProtos) // output:[h2 http/1.1]

    return cli, nil
}

// 哪些error可以确定服务端一定没有收到可以重新发起请求
func ErrorCanRetry(err error) bool {
    if opErr, ok := err.(*net.OpError); ok {
        if sysErr, ok := opErr.Err.(*os.SyscallError); ok {
            errno, ok := sysErr.Err.(syscall.Errno)
            if ok {
                /*errno == syscall.ECONNABORTED ||*/
                if errno == syscall.ECONNRESET || errno == syscall.ECONNREFUSED || errno == syscall.ENETUNREACH || errno == syscall.EPIPE {
                    return true
                }
            }
        }
    }
    return false
}
