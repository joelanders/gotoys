package main
import (
    "fmt"
    "encoding/hex"
    "encoding/base64"
)

func main() {
    b = make([]byte, 1024)
    shex := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
    n, err := hex.Decode(b, []byte(shex))
    if err != nil {
        panic(err)
    }
    s64 := base64.StdEncoding.EncodeToString(b[:n])
    fmt.Println(b[:n])
    fmt.Println(s64)
}
