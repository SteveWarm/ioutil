package smartnet

import (
    "io"
    "net"
)

var _ net.Listener = (*ChanListener)(nil)

type ChanListener struct {
    addr net.Addr
    ch   <-chan net.Conn
}

func NewChanListener(addr net.Addr, ch <-chan net.Conn) *ChanListener {
    return &ChanListener{addr: addr, ch: ch}
}

func (this *ChanListener) Accept() (net.Conn, error) {
    c, ok := <-this.ch
    if ok {
        return c, nil
    } else {
        return nil, io.EOF
    }
}

func (this *ChanListener) Close() error {
    return nil
}

// Addr returns the listener's network address.
func (this *ChanListener) Addr() net.Addr {
    return this.addr
}
