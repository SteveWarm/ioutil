package fileline

import (
    "fmt"
    "strings"
    "testing"
    "time"
)

// build:noinline
func abc() {
    fmt.Println(FileLine())
    time.Sleep(10 * time.Millisecond)
}

// build:noinline
func efg() {
    abc()
}

func TestFileLine(t *testing.T) {
    for i := 0; i < 3; i++ {
        go func() {
            abc()
        }()

    }

    for i := 0; i < 6; i++ {
        go func() {
            efg()
        }()
    }

    fl := NewCachedFileLine()
    a := fl("a")
    b := fl("a")
    if !strings.HasPrefix(a, "fileline_test.go") || !strings.HasPrefix(b, "fileline_test.go") {
        t.Fatal(a, b)
    }
    if a != b {
        t.Fatal(a, b)
    }
    t.Log(b)

    c := fl("c")
    if c == a || c == b {
        t.Fatal(a, b, c)
    }

    if !strings.HasPrefix(c, "fileline_test.go") {
        t.Fatal(c)
    }
    t.Log(c)

    time.Sleep(100 * time.Second)
}
