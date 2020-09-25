package smartnet

import (
    "fmt"
    "net"
    "time"
)

var _ net.Conn = (*Conn)(nil)

type StackBuffer = [10]byte

type Conn struct {
    c        net.Conn
    buff     StackBuffer
    buffSize int
    cur      int
}

func NewConn(conn net.Conn, buff StackBuffer, buffSize int) *Conn {
    if buffSize > len(StackBuffer{}) {
        panic(fmt.Sprintf("NewConn: buffSize: %d > > len(StackBuffer{}): %d", buffSize, len(StackBuffer{})))
    }
    return &Conn{c: conn, buff: buff, buffSize: buffSize}
}

func (this *Conn) Read(b []byte) (n int, err error) {
    if this.cur < this.buffSize {
        n := copy(b, this.buff[this.cur:this.buffSize])
        if n > 0 {
            this.cur += n
            return n, nil
        } else {
            panic(fmt.Sprintln("[BUG] SmartConn.Read n:", n, "<= 0", "cur:", this.cur))
            return n, nil // BUG
        }
    }
    return this.c.Read(b)
}

func (this *Conn) Write(b []byte) (n int, err error) {
    return this.c.Write(b)
}

func (this *Conn) Close() error {
    return this.c.Close()
}

func (this *Conn) LocalAddr() net.Addr {
    return this.c.LocalAddr()
}

func (this *Conn) RemoteAddr() net.Addr {
    return this.c.RemoteAddr()
}

func (this *Conn) SetDeadline(t time.Time) error {
    return this.c.SetDeadline(t)
}

func (this *Conn) SetReadDeadline(t time.Time) error {
    return this.c.SetReadDeadline(t)
}

func (this *Conn) SetWriteDeadline(t time.Time) error {
    return this.c.SetWriteDeadline(t)
}
