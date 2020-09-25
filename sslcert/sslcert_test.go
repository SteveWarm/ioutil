package sslcert

import (
    "testing"
)

func TestGenerateCert(t *testing.T) {
    key, cert, err := GenerateCert()
    if err != nil {
        t.Fatal(err)
    }
    t.Log("\n" + string(key))
    t.Log("\n" + string(cert))
}
