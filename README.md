# ioutil
golang io工具集,从原生函数上提炼通用的io功能，方便成组合更强大功能。
目前支持Find文件查找和Grep文本搜索。

# ini

ini格式解析和序列化，直接访问字典方便通过前缀实现类数组配置

## 样例

# Find
指定初始目录，扫描返回所有符合条件的文件和目录

## 样例
```cassandraql
	frr, drr, err := Find(Config{
		Dir:        "/tmp",
		AppendFile: true,
		AppendDir:  true,
		MinDepth:   2,
		MaxDepth:   3})
```

# Grep
从文本或文件中按行查找符合匹配的字符串

## 样例
```cassandraql
	r := strings.NewReader("a=c\nadmin_monitor_listen_port = 12345\nadmin_monitor_listen_addr=addr")
	lines, err := GrepFromReader(r, "^(admin_monitor_\\S+)[^=]*=\\s*(\\S+)$")
```

```cassandraql
	lines, err := GrepFromFile("/tmp/testgrep", "^(admin_monitor_\\S+)[^=]*=\\s*(\\S+)$")
```


# FileLine
获取的代码文件名和行号，格式为xxxx.go:1234

## 样例

```cassandraql
// global
var _fileline_ = ioutil.NewCachedFileLine()

// use 高频使用用缓存
log.Warn("some thing wrong",_fileline_("1.push.jjj.com"))
log.Warn("some thing wrong",_fileline_("2.push.jjj.com"))

// 低频使用
log.Warn("some thing wrong",ioutil.FileLine())
```


# weixinmp

微信公众平台开发接口封装

# watch

计时器；利用String反射自动调用计时。

## 样例

```
w:=Watch(time.Now())
defer log.Info("us", w)
do something ...
```

# syscmd

调用系统命令封装，自动查找第一个启动应用

# eth

读取本机ip信息

# cgroup

cgroup尝试，通过cgroup限制进程的内存cpu等使用

# httpclient

支持http http2 ssl认证 ssl双向认证 Cookies的HttpClient封装