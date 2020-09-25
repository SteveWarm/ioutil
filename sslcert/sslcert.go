// 生成SSL相关证书
package sslcert

//
// 参考资料
// 推荐书籍《图解密码技术(日)结城浩(著)》
// PKCS 15 个标准 https://www.cnblogs.com/jtlgb/p/6762050.html
// https://blog.csdn.net/fyxichen/article/details/53010255
// http://www.hydrogen18.com/blog/your-own-pki-tls-golang.html
//

import (
    "bytes"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "math/big"
    rd "math/rand"
    "time"
)

func GenerateCert() ([]byte, []byte, error) {
    // 生成ca.key
    pk, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, nil, err
    }

    // 生成ca.crt
    serialNumber, err := rand.Int(rand.Reader, big.NewInt(rd.Int63()))
    if err != nil {
        return nil, nil, err
    }
    template := &x509.Certificate{
        SerialNumber: serialNumber,
        Subject: pkix.Name{
            Organization: []string{"woodada"},
            CommonName:   "woodada", // Will be checked by the server
            // 不止两个参数还有一堆要填的
        },
        NotBefore:             time.Now(),
        NotAfter:              time.Now().Add(30 * 24 * time.Hour),
        KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
        BasicConstraintsValid: true,
    }

    header := make(map[string]string)
    pkData := x509.MarshalPKCS1PrivateKey(pk)                    // RSA PRIVATE KEY
    keyData, err := pemEncode("RSA PRIVATE KEY", pkData, header) // /tmp/rsa.key
    if err != nil {
        return nil, nil, err
    }
    // ed25519 rsa 比较省心的是 Marshal函数注释里备注了Type名 不用再去查资料
    // x509.MarshalPKIXPublicKey() // PUBLIC KEY
    // 自签名
    certData, err := x509.CreateCertificate(rand.Reader, template, template, pk.Public(), pk)
    if err != nil {
        return nil, nil, err
    }
    certData, err = pemEncode("CERTIFICATE", certData, header) // /tmp/rsa.crt
    if err != nil {
        return nil, nil, err
    }
    return keyData, certData, nil
}

func pemEncode(typ string, data []byte, header map[string]string) ([]byte, error) {
    w := &bytes.Buffer{}
    err := pem.Encode(w, &pem.Block{Bytes: data, Type: typ, Headers: header})
    if err != nil {
        return nil, err
    }
    return w.Bytes(), nil
}
