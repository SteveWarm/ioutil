/*
 网络相关的封装放本文件
*/
package eth

import (
    "net"
)

type Eth struct {
    Name         string
    IP4Addr      []string
    IP4IntAddr   []uint32
    HardwareAddr string
    Flags        net.Flags
}

// 获取本地所有活跃的网卡信息
func LiveEths() ([]Eth, error) {
    infs, err := net.Interfaces()
    if err != nil {
        return nil, err
    }

    var eths []Eth
    for _, inf := range infs {
        if (inf.Flags & net.FlagUp) != 0 {
            eth := Eth{}
            eth.Flags = inf.Flags
            eth.Name = inf.Name
            eth.HardwareAddr = inf.HardwareAddr.String()

            arr, err := inf.Addrs()
            if err != nil {
                continue
            }

            for _, a := range arr {
                if ip, ok := a.(*net.IPNet); ok {
                    if ip.IP.To4() != nil {
                        if ip.IP.String() != "" {
                            eth.IP4Addr = append(eth.IP4Addr, ip.IP.String())
                            eth.IP4IntAddr = append(eth.IP4IntAddr, IPv4StrToUint32(ip.IP.String()))
                        }
                    }
                }
            }

            if len(eth.IP4Addr) > 0 {
                eths = append(eths, eth)
            }
        }
    }

    return eths, nil
}
