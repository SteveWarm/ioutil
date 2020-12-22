// 微信开放平台开发接口封装
package weixinmp

import (
    "bytes"
    "encoding/json"
    "errors"

    "net/http"
    "sync"
    "time"
)

type QQ WeiXin

func NewQQ(config WeiXinConfig) *QQ {
    config.Jscode2sessionURL = firstAvailable(config.Jscode2sessionURL, "https://api.q.qq.com/sns/jscode2session")
    config.GetAccessTokenURL = firstAvailable(config.GetAccessTokenURL, "https://api.q.qq.com/api/getToken")
    config.MsgSecCheckURL = firstAvailable(config.MsgSecCheckURL, "https://api.q.qq.com/api/json/security/MsgSecCheck")
    return (*QQ)(NewWeiXin(config))
}

const (
    code2session_url = "https://api.weixin.qq.com/sns/jscode2session"
    access_token_url = "https://api.weixin.qq.com/cgi-bin/token"
    sec_check_url    = "https://api.weixin.qq.com/wxa/msg_sec_check"
)

type WeiXinConfig struct {
    AppID             string
    Secret            string
    Jscode2sessionURL string
    GetAccessTokenURL string
    MsgSecCheckURL    string
}

type WeiXin struct {
    config             WeiXinConfig
    accessTokenInfo    AccessTokenInfo
    accessTokenInfoMux sync.RWMutex
}

func firstAvailable(ss ...string) string {
    for _, s := range ss {
        if s != "" {
            return s
        }
    }
    return ""
}

func NewWeiXin(config WeiXinConfig) *WeiXin {
    config.Jscode2sessionURL = firstAvailable(config.Jscode2sessionURL, code2session_url)
    config.GetAccessTokenURL = firstAvailable(config.GetAccessTokenURL, access_token_url)
    config.MsgSecCheckURL = firstAvailable(config.MsgSecCheckURL, sec_check_url)
    w := &WeiXin{config: config}
    go w.update_access_token_runable()
    return w
}

func (this *WeiXin) GetAccessTokenInfo() AccessTokenInfo {
    this.accessTokenInfoMux.RLock()
    defer this.accessTokenInfoMux.RUnlock()
    return this.accessTokenInfo
}

// 内容安全：校验文本内容是否合规
// errCode	number	错误码
// errMsg	string	错误信息
// 小于0系统内部错误 大于0业务错误 0 成功
func (this *WeiXin) MsgSecCheck(content string) (int64, string) {
    type Req struct {
        Content string `json:"content"`
    }

    data, err := json.Marshal(Req{Content: content})
    if err != nil {
        return -1, err.Error()
    }

    r := bytes.NewReader(data)
    req, err := http.NewRequest("POST", this.config.MsgSecCheckURL, r)
    if err != nil {
        return -1, err.Error()
    }
    req.Header.Set("Content-type", "application/json")
    q := req.URL.Query()
    q.Add("appid", this.config.AppID)
    q.Add("access_token", this.accessTokenInfo.AccessToken)
    req.URL.RawQuery = q.Encode()

    rsp, err := http.DefaultClient.Do(req)
    defer func() {
        if rsp != nil && rsp.Body != nil {
            rsp.Body.Close()
        }
    }()

    if err != nil {
        return -2, err.Error()
    }

    type Rsp struct {
        ErrCode int64  `json:"errCode"` // errCode	number	错误码
        ErrMsg  string `json:"errMsg"`  // errMsg	string	错误信息
    }

    if rsp.StatusCode == 200 {
        v := &Rsp{}
        err := json.NewDecoder(rsp.Body).Decode(v)
        if err == nil {
            if v.ErrCode == 0 {
                logger.Info("MsgSecCheck succ! content:", content, "result:", * v)
                return v.ErrCode, v.ErrMsg
            } else {
                logger.Warn("MsgSecCheck fail! content:", content, "result:", *v)
                return v.ErrCode, v.ErrMsg
            }
        } else {
            logger.Error("MsgSecCheck Decoder err:", err)
            return -1, err.Error()
        }
    } else {
        logger.Warn("MsgSecCheck Request fail! Status:", rsp.Status)
        return -1, rsp.Status
    }
}

func (this *WeiXin) update_access_token_runable() {
    for {
        sleepSecond := 1 * time.Second
        func() {
            defer func() {
                if err := recover(); err != nil {
                    logger.Error("recover", err)
                }
            }()

            // 过期更新
            info, err := this.getAccessToken()
            if err != nil {
                logger.Error("getAccessToken err:", err)
                return
            }

            if info.ErrCode == 0 {
                logger.Info("getAccessToken succ! resp:", *info)
                sleepSecond = time.Duration(info.ExpiresIn-600) * time.Second // 提前10分钟更新
                this.accessTokenInfoMux.Lock()
                defer this.accessTokenInfoMux.Unlock()
                this.accessTokenInfo = *info
            } else {
                sleepSecond = 1 * time.Second
                logger.Warn("getAccessToken fail! resp:", info)
            }
        }()
        time.Sleep(sleepSecond)
    }
}

type AccessTokenInfo struct {
    AccessToken string `json:"access_token"`
    ExpiresIn   int64  `json:"expires_in"` // 凭证有效时间，单位：秒。目前是7200秒之内的值。
    ErrMsg      string `json:"errmsg"`     // 错误信息
    // ErrCode 错误码
    // -1	系统繁忙，此时请开发者稍候再试
    // 0	请求成功
    // 40001	AppSecret 错误或者 AppSecret 不属于这个小程序，请开发者确认 AppSecret 的正确性
    // 40002	请确保 grant_type 字段值为 client_credential
    // 40013	不合法的 AppID，请开发者检查 AppID 的正确性，避免异常字符，注意大小写
    ErrCode int64 `json:"errcode"` // 错误码
}

func (me AccessTokenInfo) ToJson() string {
    data, _ := json.Marshal(me)
    return string(data)
}

/**
* https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
* -1 系统繁忙，此时请开发者稍候再试
* 0 请求成功
* 40029 code 无效
* 45011 频率限制，每个用户每分钟100次
 */
func (this *WeiXin) getAccessToken() (*AccessTokenInfo, error) {
    req, err := http.NewRequest("GET", this.config.GetAccessTokenURL, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-type", "application/json")
    q := req.URL.Query()
    q.Add("appid", this.config.AppID)
    q.Add("secret", this.config.Secret)
    q.Add("grant_type", "client_credential")
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
        v := &AccessTokenInfo{}
        err := json.NewDecoder(rsp.Body).Decode(v)
        return v, err
    } else {
        return nil, errors.New(rsp.Status)
    }
}
