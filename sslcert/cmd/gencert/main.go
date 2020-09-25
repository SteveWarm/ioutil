package main

import (
    "fmt"
    "io/ioutil"
    "os"

    "github.com/woodada/ioutil/sslcert"
)

func main() {
    var keyFile, certFile string
    if len(os.Args) == 3 {
        keyFile = os.Args[1]
        certFile = os.Args[2]
    } else if len(os.Args) == 1 {
        // ok
    } else {
        fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "key_output_file_path", "cert_output_file_path")
        os.Exit(1)
    }

    keyData, certData, err := sslcert.GenerateCert()
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(1)
    }

    if len(keyFile) > 0 && len(certFile) > 0 {
        e1 := ioutil.WriteFile(keyFile, keyData, os.FileMode(0644))
        e2 := ioutil.WriteFile(certFile, certData, os.FileMode(0644))
        if e1 != nil || e2 != nil {
            fmt.Fprintln(os.Stderr, e1, e2)
            os.Exit(1)
        }
    } else {
        fmt.Println(string(keyData))
        fmt.Println(string(certData))
    }

}
