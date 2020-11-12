package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "os/signal"
    "path/filepath"
    "time"
)

func main() {
    rootdir := "/sys/fs/cgroup"
    pid := os.Getpid()
    memdir := filepath.Join(rootdir, "memory", fmt.Sprint(pid))
    err := os.MkdirAll(memdir, 0644)
    assertError(err, "MkdirAll")
    err = ioutil.WriteFile(filepath.Join(memdir, "memory.limit_in_bytes"), []byte("10M"), 0644)
    assertError(err, "WriteFile memory.limit_in_bytes")
    err = ioutil.WriteFile(filepath.Join(memdir, "memory.swappiness"), []byte("0"), 0644)
    assertError(err, "WriteFile memory.swappiness")

    // 从第二个参数开始作为托管进程的启动参数
    var args []string
    for i := 0; i+1 < len(os.Args); i++ {
        args = append(args, os.Args[1+i])
    }

    go startCmd(args, memdir)

    ch := make(chan os.Signal)
    signal.Notify(ch, os.Kill, os.Interrupt)
    fmt.Println(<-ch)
    err = os.RemoveAll(memdir)
    if err != nil {
        fmt.Println(err)
    }
}

func assertError(err error, tag ...string) {
    if err != nil {
        fmt.Fprintln(os.Stderr, tag, err)
        os.Exit(1)
    }
}

func startCmd(args []string, memdir string) {
    for {
        func() {
            defer func() {
                if err := recover(); err != nil {
                    fmt.Println(err)
                }
            }()
            cmd := exec.Cmd{
                Path: args[0],
                Args: args,
                Env:  os.Environ(),
            }

            cmd.Stdin = os.Stdin
            cmd.Stderr = os.Stderr
            cmd.Stdout = os.Stdout

            // start app
            if err := cmd.Start(); err != nil {
                fmt.Println(err)
                return
            }

            // set cgroup procs id
            err := ioutil.WriteFile(filepath.Join(memdir, "cgroup.procs"), []byte(fmt.Sprint(cmd.Process.Pid)), 0644)
            if err != nil {
                fmt.Println("add pid", cmd.Process.Pid, "to file cgroup.procs", err)
            }
            if err := cmd.Wait(); err != nil {
                fmt.Println("cmd return with error:", err)
            }

            //status := cmd.ProcessState.Sys().(syscall.WaitStatus)
            //
            ////options := ExitStatus{
            ////    Code: status.ExitStatus(),
            ////}
            //
            //if status.Signaled() {
            //    options.Signal = status.Signal()
            //}
            cmd.Process.Kill()
            fmt.Println(cmd.ProcessState.String())
        }()

        time.Sleep(5 * time.Second)
    }
}
