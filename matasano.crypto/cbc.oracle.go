package main
import (
    "fmt"
    "math/rand"
    crand "crypto/rand" //use both, for fun
    "crypto/cipher" //use both, for fun
    "crypto/aes" //use both, for fun
    "bytes"
    "time"
    "io"
    "./ecb"
    "./cbc"
    "encoding/hex"
)

var r *rand.Rand

func main() {
    r = rand.New(rand.NewSource(time.Now().UnixNano()))
    encrypter := RandEncrypter(UnknownAESCipher())
    chosenPlaintext := bytes.Repeat([]byte("a"), 256)
    cipherText := make([]byte, 256)
    encrypter.CryptBlocks(cipherText, chosenPlaintext)
    fmt.Println(hex.EncodeToString(chosenPlaintext))
    fmt.Println(hex.EncodeToString(cipherText))
    dupes := ecb.CheckDupeBlock(cipherText, 16)
    if dupes {
        fmt.Println("dupe blocks, so it's ecb")
    } else {
        fmt.Println("no dupe blocks, so not ecb")
    }
}

func RandEncrypter(block cipher.Block) cipher.BlockMode {
    if r.Int()%2 == 0 {
        return ecb.NewECBEncrypter(block)
    } else {
        iv := make([]byte, 16)
        io.ReadFull(crand.Reader, iv)
        return cbc.NewCBCEncrypter(block, iv)
    }
}

func UnknownAESCipher() cipher.Block {
    key := make([]byte, 16)
    io.ReadFull(crand.Reader, key)
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }
    return block
}
