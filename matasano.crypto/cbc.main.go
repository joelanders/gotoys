package main
import (
    "fmt"
    "./cbc"
    "./xor"
    "crypto/aes"
    "encoding/base64"
)

func main() {
    bs := xor.BytesFromFile("10.txt")
    fmt.Println(len(bs))

    key := []byte("YELLOW SUBMARINE")
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }

    iv := make([]byte, 16)
    fmt.Println(iv)

    dst := make([]byte, len(bs))
    decrypter := cbc.NewCBCDecrypter(block, iv)
    decrypter.CryptBlocks(dst, bs)
    fmt.Println(string(dst))

    dst2 := make([]byte, len(bs))
    encrypter := cbc.NewCBCEncrypter(block, iv)
    encrypter.CryptBlocks(dst2, dst)
    fmt.Println(base64.StdEncoding.EncodeToString(dst2))
}

