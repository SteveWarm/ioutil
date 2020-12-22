package weixinmp

import (
    "encoding/json"
    "errors"
    "net/http"
)

type AuthInfo struct {
    OpenID     string `json:""openid"`     // 用户唯一标识
    SessionKey string `json:"session_key"` // 会话密钥
    UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符，在满足 UnionID 下发条件的情况下会返回，详见 UnionID 机制说明。
    ErrCode    int64  `json:"errcode"`     // 错误码
    ErrMsg     string `json:"errmsg"`      // 错误信息
}

func (this *AuthInfo) JsonStr() string {
    d, _ := json.Marshal(this)
    return string(d)
}

/**
* 根据登陆信息获取openid
* https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
* -1 系统繁忙，此时请开发者稍候再试
* 0 请求成功
* 40029 code 无效
* 45011 频率限制，每个用户每分钟100次
 */
func JsCode2Session(url, appid, secret, jscode string) (*AuthInfo, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    q := req.URL.Query()
    q.Add("appid", appid)
    q.Add("secret", secret)
    q.Add("js_code", jscode)
    q.Add("grant_type", "authorization_code")
    req.URL.RawQuery = q.Encode()

    rsp, err := http.DefaultClient.Do(req)
    defer func() {
        if rsp != nil && rsp.Body != nil {
            rsp.Body.Close()
        }
    }()

    if err != nil {
        return nil, err
    }

    if rsp.StatusCode == 200 {
        v := &AuthInfo{}
        err := json.NewDecoder(rsp.Body).Decode(v)
        return v, err
    } else {
        return nil, errors.New(rsp.Status)
    }
}
