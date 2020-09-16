# ioutil
golang io工具集,从原生函数上提炼通用的io功能，方便成组合更强大功能。
目前支持Find文件查找和Grep文本搜索。

# Find
指定初始目录，扫描返回所有符合条件的文件和目录

## 定义
```cassandraql
func Find(config Config) (files []string, dirs []string, err error) 

type Config struct {
	Dir        string   // 需要扫描的起始目录
	MinDepth   int      // 最小深度 默认0每层都记录
	MaxDepth   int      // 最大深度 超过本深度不探测
	AppendFile bool     // 是否返回文件列表 全路径
	AppendDir  bool     // 是否返回目录列表 全路径
	RegexMatch []string // 路径正则匹配，符合匹配的才返回和向下探测
	_regexList []*regexp.Regexp // 内部使用，RegexMatch编译呴缓存
}

```


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

## 定义
```cassandraql
func GrepFromReader(r io.Reader, regstr string) (lines []string, err error)
func GrepFromFile(filestr, regstr string) (lines []string, err error) 
```
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

## 定义

```cassandraql
func NewFileLineCache() func(key string) string
func FileLine() string 
```
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