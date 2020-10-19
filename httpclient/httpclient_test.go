package httpclient

import (
    "io/ioutil"
    "net/http"
    "testing"
)

func TestHttpClient(t *testing.T) {
    opt := Options{}

    c, err := NewHttpClient(opt)
    if err != nil {
        t.Fatal(err)
    }

    testcases := []string{"https://baidu.com", "https://whois.pconline.com.cn/ip.jsp", "https://zhihu.com"}

    for _, url := range testcases {
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            t.Fatal(err)
        }
        rsp, err := c.Do(req)
        if err != nil {
            t.Fatal(err)
        }

        defer rsp.Body.Close()

        t.Log(rsp.Status)
        data, err := ioutil.ReadAll(rsp.Body)
        if err != nil {
            t.Fatal(err)
        }
        // t.Log(string(data))
        unuse(data)
        if rsp.StatusCode != 200 {
            t.Fail()
        }
    }
}

func unuse(interface{}) {

}
