package eth

import (
    "testing"
)

func TestPubblicAddr(t *testing.T) {
    pa, err := PubblicAddr("")
    if err != nil {
        t.Error("请联系作者更新接口实现", err)
        t.FailNow()
    }
    t.Log(pa)

    eths, err := LiveEths()
    if err != nil {
        t.Error(err)
        t.FailNow()
    }

    for _, eth := range eths {
        if len(eth.IP4Addr) > 0 {
            pa, err := PubblicAddr(eth.IP4Addr[0])
            t.Log(eth.Name, pa, err)
        }
    }
}
