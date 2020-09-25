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