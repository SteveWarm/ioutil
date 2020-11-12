package eth

import (
    "context"
    "encoding/json"
    "github.com/axgle/mahonia"
    "io/ioutil"
    "net"
    "net/http"
    "strings"
)

type NetAddr struct {
    Addr        string `json:"addr"`        // 地址
    City        string `json:"city"`        // 城市
    CityCode    string `json:"cityCode"`    // 城市邮政编码
    Err         string `json:"err"`         // 如果有错误的话提示
    IP          string `json:"ip"`          // 公网ip
    Pro         string `json:"pro"`         // 省
    ProCode     string `json:"proCode"`     // 省邮政编码
    Region      string `json:"region"`      //
    RegionCode  string `json:"regionCode"`  //
    RegionNames string `json:"regionNames"` //
}

// 公网地址
// localIP 指定本机走那个网卡地址，可以为空
func PubblicAddr(localIP string) (NetAddr, error) {
    c := http.Client{
        Transport: &http.Transport{
            DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
                l, e := net.ResolveTCPAddr(network, localIP+":0")
                if e != nil {
                    return nil, e
                }

                r, e := net.ResolveTCPAddr(network, addr)
                if e != nil {
                    return nil, e
                }
                return net.DialTCP(network, l, r)
            },
        },
    }

    // http://whois.pconline.com.cn/ipJson.jsp
    rsp, err := c.Get("http://whois.pconline.com.cn/ipJson.jsp")
    if err != nil {
        return NetAddr{}, err
    }

    defer rsp.Body.Close()

    contentType := rsp.Header.Get("Content-Type")
    charset := "GBK"
    for _, v := range strings.Split(contentType, ";") {
        v = strings.TrimSpace(v)
        if strings.HasPrefix(v, "charset=") {
            charset = strings.TrimPrefix(v, "charset=")
            break
        }
    }

    data, err := ioutil.ReadAll(rsp.Body)
    if err != nil {
        return NetAddr{}, err
    }

    addr := mahonia.NewDecoder(charset).ConvertString(string(data))
    addr = strings.TrimSpace(addr)
    addr = strings.TrimPrefix(addr, "if(window.IPCallBack) {IPCallBack(")
    addr = strings.TrimSuffix(addr, ");}")

    var pa NetAddr

    err = json.NewDecoder(strings.NewReader(addr)).Decode(&pa)
    if err != nil {
        return pa, err
    }

    return pa, nil
}
