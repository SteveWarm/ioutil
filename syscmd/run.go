package syscmd

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
)

type Error int

func (me Error) Error() string {
    switch me {
    case ErrCmdNotFound:
        return "ErrCmdNotFound"
    default:
        return "unknown:" + strconv.Itoa(int(me))
    }
}

const (
    ErrCmdNotFound Error = 1
)

// 运行一段命令
// Args: cmdline
// Return: stat stdout errout error
func RunStr(cmdline string) (string, string, string, error) {
    arr := strings.Split(cmdline, " ")
    var args []string
    for _, a := range arr {
        if len(a) <= 0 {
            continue
        }
        args = append(args, a)
    }

    return RunArgs(args)
}

// 运行一段命令，将入参拆解
// Args: cmdline
// Return: stat stdout errout error
func RunArgs(args []string) (string, string, string, error) {
    fileInfo, err := os.Stat(args[0])
    if err == nil {
        if fileInfo.IsDir() {
            return "", "", "", fmt.Errorf("IsDirErr:%s", args[0])
        }
    } else {
        rawPath, err := Which(args[0])
        if err != nil {
            return "", "", "", err
        }
        args[0] = rawPath
    }

    // 始终用绝对路径执行
    args[0], err = filepath.Abs(args[0])
    if err != nil {
        return "", "", "", err
    }

    // cmdstr := strings.Join(args, " ")
    data, _ := json.Marshal(args)
    cmdjson := string(data)
    //  log("[CMD-STR]", cmdstr)
    log("[CMD-JSON]", cmdjson)
    if NoRunEnabled() {
        return "", "", "", nil
    }
    stdout := &strings.Builder{}
    errout := &strings.Builder{}
    cmd := exec.Cmd{
        Path: args[0],
        Args: args,
        Env:  os.Environ(),
    }

    cmd.Stdin = os.Stdin
    cmd.Stderr = errout
    cmd.Stdout = stdout

    // start app
    if err := cmd.Start(); err != nil {
        log(err)
        return "", "", "", err
    }

    err = cmd.Wait()

    //status := cmd.ProcessState.Sys().(syscall.WaitStatus)
    //
    ////options := ExitStatus{
    ////    Code: status.ExitStatus(),
    ////}
    //
    //if status.Signaled() {
    //    options.Signal = status.Signal()
    //}
    //err := cmd.Process.Kill()
    return cmd.ProcessState.String(), stdout.String(), errout.String(), err
}
