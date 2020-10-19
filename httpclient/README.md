支持http http2 ssl认证 ssl双向认证 Cookies的HttpClient封装

```
NewHttpClient(option Options) (*http.Client, error) 
```

```
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
```