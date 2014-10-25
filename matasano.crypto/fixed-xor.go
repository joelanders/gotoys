package main
import (
    "fmt"
    "encoding/hex"
)

func main() {
    ahex := "1c0111001f010100061a024b53535009181c"
    bhex := "686974207468652062756c6c277320657965"

    a, err := hex.DecodeString(ahex)
    if err != nil {
        panic(err)
    }
    b, err := hex.DecodeString(bhex)
    if err != nil {
        panic(err)
    }

    fmt.Println(a)
    fmt.Println(b)

    fmt.Println(string(xorSlices(a,b)))
}

func xorSlices(a, b []byte) []byte {
    if len(a) != len(b) {
        panic("unequal lengths")
    }
    
    c := make([]byte, len(a))

    for i, _ := range(a) {
        c[i] = a[i] ^ b[i]
    }
    return c
}
