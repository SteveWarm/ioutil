package syscmd

import (
    "strings"
)

// 查找某个命令的路径
// 不同平台用不同实现
// TODO 目前只支持linux mac
func Which(s string) (string, error) {
    _, stdout, _, err := RunArgs([]string{"/usr/bin/which", s})
    rawPath := strings.TrimSpace(stdout)
    if err == nil && len(rawPath) <= 0 {
        return "", ErrCmdNotFound
    }
    return strings.TrimSpace(stdout), err
}
