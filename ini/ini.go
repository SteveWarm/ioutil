// ini格式解析和序列化
package ini

import (
    "bufio"
    "bytes"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "regexp"
    "strconv"
    "strings"
    "unicode"
)

type Dict map[string]map[string]string

var (
    regDoubleQuote = regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*\"([^\"]*)\"$")
    regSingleQuote = regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*'([^']*)'$")
    regNoQuote     = regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*([^#;]+)")
    regNoValue     = regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*([#;].*)?")
)

func Parse(data []byte) (Dict, error) {
    reader := bytes.NewReader(data)
    return parse(reader, "data")
}

func Load(filename string) (Dict, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    return parse(file, filename)
}

func parse(r io.Reader, filename string) (dict Dict, err error) {
    reader := bufio.NewReader(r)
    lineno := 0
    section := ""
    dict[section] = make(map[string]string)

    for err == nil {
        l, _, err := reader.ReadLine()
        if err != nil {
            break
        }
        lineno++
        if len(l) == 0 {
            continue
        }
        line := strings.TrimFunc(string(l), unicode.IsSpace)

        for len(line) > 0 && line[len(line)-1] == '\\' {
            line = line[:len(line)-1]
            l, _, err := reader.ReadLine()
            if err != nil {
                return nil, err
            }
            line += strings.TrimFunc(string(l), unicode.IsSpace)
        }

        section, err = dict.parseLine(section, line)
        if err != nil {
            return nil, errors.New(
                err.Error() + fmt.Sprintf("'%s:%d'.", filename, lineno))
        }
    }
    return
}

func Write(filename string, dict *Dict) error {
    buffer := dict.format()
    return ioutil.WriteFile(filename, buffer.Bytes(), 0644)
}

func (dict Dict) parseLine(section, line string) (string, error) {
    if len(line) <= 0 {
        return section, nil
    }

    // commets
    if line[0] == '#' || line[0] == ';' {
        return section, nil
    }

    // section name
    if line[0] == '[' && line[len(line)-1] == ']' {
        section := strings.TrimFunc(line[1:len(line)-1], unicode.IsSpace)
        section = strings.ToLower(section)
        dict[section] = make(map[string]string)
        return section, nil
    }

    // key = value
    if m := regDoubleQuote.FindAllStringSubmatch(line, 1); m != nil {
        dict.add(section, m[0][1], m[0][2])
        return section, nil
    } else if m = regSingleQuote.FindAllStringSubmatch(line, 1); m != nil {
        dict.add(section, m[0][1], m[0][2])
        return section, nil
    } else if m = regNoQuote.FindAllStringSubmatch(line, 1); m != nil {
        dict.add(section, m[0][1], strings.TrimFunc(m[0][2], unicode.IsSpace))
        return section, nil
    } else if m = regNoValue.FindAllStringSubmatch(line, 1); m != nil {
        dict.add(section, m[0][1], "")
        return section, nil
    }

    return section, errors.New("iniparser: syntax error at ")
}

func (dict Dict) add(section, key, value string) {
    key = strings.ToLower(key)
    dict[section][key] = value
}

func (dict Dict) GetBool(section, key string) (bool, bool) {
    sec, ok := dict[section]
    if !ok {
        return false, false
    }
    value, ok := sec[key]
    if !ok {
        return false, false
    }
    v := value[0]
    if v == 'y' || v == 'Y' || v == '1' || v == 't' || v == 'T' {
        return true, true
    }
    if v == 'n' || v == 'N' || v == '0' || v == 'f' || v == 'F' {
        return false, true
    }
    return false, false
}

func (dict Dict) SetBool(section, key string, value bool) {
    dict.SetString(section, key, strconv.FormatBool(value))
}

func (dict Dict) GetString(section, key string) (string, bool) {
    sec, ok := dict[section]
    if !ok {
        return "", false
    }
    value, ok := sec[key]
    if !ok {
        return "", false
    }
    return value, true
}

func (dict Dict) SetString(section, key, value string) {
    _, ok := dict[section]
    if !ok {
        dict[section] = make(map[string]string)
    }
    dict[section][key] = value
}

func (dict Dict) GetInt(section, key string) (int, bool) {
    sec, ok := dict[section]
    if !ok {
        return 0, false
    }
    value, ok := sec[key]
    if !ok {
        return 0, false
    }
    i, err := strconv.Atoi(value)
    if err != nil {
        return 0, false
    }
    return i, true
}

func (dict Dict) SetInt(section, key string, value int) {
    dict.SetString(section, key, strconv.FormatInt(int64(value), 10))
}

func (dict Dict) GetDouble(section, key string) (float64, bool) {
    sec, ok := dict[section]
    if !ok {
        return 0, false
    }
    value, ok := sec[key]
    if !ok {
        return 0, false
    }
    d, err := strconv.ParseFloat(value, 64)
    if err != nil {
        return 0, false
    }
    return d, true
}

func (dict Dict) SetDouble(section, key string, value float64) {
    dict.SetString(section, key, strconv.FormatFloat(value, 'f', -1, 64))
}

func (dict Dict) Delete(section, key string) {
    _, ok := dict[section]
    if !ok {
        return
    }
    delete(dict[section], key)
}

func (dict Dict) GetSections() []string {
    size := len(dict)
    sections := make([]string, size)
    i := 0
    for section, _ := range dict {
        sections[i] = section
        i++
    }
    return sections
}

func (dict Dict) GetKeys(section string) []string {
    sec, ok := dict[section]
    if !ok {
        return nil
    }
    size := len(sec)
    if size <= 0 {
        return nil
    }

    keys := make([]string, size)
    i := 0
    for k, _ := range sec {
        keys[i] = k
        i++
    }
    return keys
}

func (dict Dict) String() string {
    return (*dict.format()).String()
}

func (dict Dict) format() *bytes.Buffer {
    var buffer bytes.Buffer
    for section, vals := range dict {
        if section != "" {
            buffer.WriteString(fmt.Sprintf("[%s]\n", section))
        }
        for key, val := range vals {
            buffer.WriteString(fmt.Sprintf("%s = %s\n", key, val))
        }
        buffer.WriteString("\n")
    }
    return &buffer
}
