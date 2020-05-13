package find

import (
	"strings"
	"testing"
)

func TestGrepFromReader(t *testing.T) {
	r := strings.NewReader("a=c\nadmin_monitor_listen_port = 12345\nadmin_monitor_listen_addr=addr")
	lines, err := GrepFromReader(r, "^(admin_monitor_\\S+)[^=]*=\\s*(\\S+)$")
	if err != nil {
		t.Fatal(err)
	}
	for i, line := range lines {
		t.Log(i, ".", line)
	}
}

func TestGrepFromFile(t *testing.T) {
	/*
		/tmp/testgrep
		a=c
		admin_monitor_listen_port = 12345
		admin_monitor_listen_addr=addr
	*/
	lines, err := GrepFromFile("/tmp/testgrep", "^(admin_monitor_\\S+)[^=]*=\\s*(\\S+)$")
	if err != nil {
		t.Fatal(err)
	}
	for i, line := range lines {
		t.Log(i, ".", line)
	}
}
