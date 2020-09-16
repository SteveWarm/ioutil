package ioutil

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	Dir        string   // 需要扫描的起始目录
	MinDepth   int      // 最小深度 默认0每层都记录
	MaxDepth   int      // 最大深度 超过本深度不探测
	AppendFile bool     // 是否返回文件列表 全路径
	AppendDir  bool     // 是否返回目录列表 全路径
	RegexMatch []string // 路径正则匹配，符合匹配的才返回和向下探测
	_regexList []*regexp.Regexp // 内部使用，RegexMatch编译呴缓存
}

func (this Config) _match_path(path string) bool {
	if this._regexList == nil {
		return true
	}

	for _, r := range this._regexList {
		if r.MatchString(path) {
			return true
		}
	}
	return false
}

func Find(config Config) (files []string, dirs []string, err error) {
	if config.RegexMatch != nil {
		for _, str := range config.RegexMatch {
			if r, e := regexp.Compile(str); e == nil {
				config._regexList = append(config._regexList, r)
			} else {
				err = e
				return
			}

		}
	}

	err = listDir(&files, &dirs, config.Dir, 1, config)
	return
}

func listDir(files *[]string, dirs *[]string, dir string, curDepth int, c Config) error {
	if !strings.HasSuffix(dir, string(os.PathSeparator)) {
		dir += string(os.PathSeparator)
	}

	rd, err := ioutil.ReadDir(dir)
	if nil != err {
		return err
	}

	if c.AppendFile && curDepth >= c.MinDepth {
		for _, fi := range rd {
			if !fi.IsDir() {
				pathstr := dir + fi.Name()
				if c._match_path(pathstr) {
					*files = append(*files, pathstr)
				}
			}
		}
	}

	for _, fi := range rd {
		if fi.IsDir() {
			if c.AppendDir && curDepth >= c.MinDepth {
				pathstr := dir + fi.Name()
				if c._match_path(pathstr) {
					*dirs = append(*dirs, pathstr)
				}
			}

			if c.MaxDepth == 0 || curDepth+1 <= c.MaxDepth {
				if err := listDir(files, dirs, dir+fi.Name(), curDepth+1, c); err != nil {
					return err
				}
			}
		}
	}

	return err
}
