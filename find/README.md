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