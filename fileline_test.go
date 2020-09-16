package ioutil

import (
    "strings"
    "testing"
)

func TestFileLine(t *testing.T) {
    t.Log(FileLine())

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

}
