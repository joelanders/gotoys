package main

import (
    "os"
    "fmt"
    "encoding/pem"
    "crypto/x509"
)

func main() {
    f, err := os.Open("cert.crt")
    if err != nil {
        panic(err)
    }

    buf := make([]byte, 8*1024)
    n, err := f.Read(buf)
    if err != nil {
        panic(err)
    }

    block, _ := pem.Decode(buf[:n])

    cert, err := x509.ParseCertificate(block.Bytes)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%v", cert)
}
