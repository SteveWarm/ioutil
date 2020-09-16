package ioutil

import (
    "runtime"
    "strconv"
    "strings"
    "sync"
)

func NewCachedFileLine() func(key string) string {
    var _lock sync.RWMutex
    var _cache map[string]string = make(map[string]string)
    return func(key string) string {
        _lock.RLock()
        v, ok := _cache[key]
        _lock.RUnlock()

        if ok {
            return v
        }

        _, file, line, ok := runtime.Caller(2) // decorate + log + public function.
        if ok {
            // Truncate file name at last file name separator.
            if index := strings.LastIndex(file, "/"); index >= 0 {
                file = file[index+1:]
            } else if index = strings.LastIndex(file, "\\"); index >= 0 {
                file = file[index+1:]
            }
        } else {
            file = "???"
            line = 1
        }

        _lock.Lock()
        v = file + ":" + strconv.Itoa(line)
        _cache[key] = v
        _lock.Unlock()
        return v
    }
}

func FileLine() string {
    _, file, line, ok := runtime.Caller(2) // decorate + log + public function.
    if ok {
        // Truncate file name at last file name separator.
        if index := strings.LastIndex(file, "/"); index >= 0 {
            file = file[index+1:]
        } else if index = strings.LastIndex(file, "\\"); index >= 0 {
            file = file[index+1:]
        }
    } else {
        file = "???"
        line = 1
    }
    return file + ":" + strconv.Itoa(line)
}
