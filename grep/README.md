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