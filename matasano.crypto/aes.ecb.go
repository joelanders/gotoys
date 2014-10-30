package main
import (
    "fmt"
    "crypto/aes"
    //"encoding/hex"
    "encoding/base64"
    "./xor" //todo move stuff
    "./ecb"
)

func main() {
    key := []byte("YELLOW SUBMARINE")
//    keyHex := "2b7e151628aed2a6abf7158809cf4f3c"
//    key, err := hex.DecodeString(keyHex)
//    if err != nil {
//        panic(err)
//    }
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }

    //this base64 decodes
    bs := xor.BytesFromFile("7.txt")
    fmt.Println(len(bs))

    //3ad77bb40d7a3660a89ecaf32466ef97
    //f5d3d58503b9699de785895a96fdbaaf
//    src1Hex := "6bc1bee22e409f96e93d7e117393172aae2d8a571e03ac9c9eb76fac45af8e51"
//    src1, err := hex.DecodeString(src1Hex)
//    if err != nil {
//        panic(err)
//    }
//    dst1 := make([]byte, 32)
//
//    encrypter := NewECBEncrypter(block)
//    encrypter.CryptBlocks(dst1, src1)
//    fmt.Println(hex.EncodeToString(dst1))

    dst := make([]byte, len(bs))
    decrypter := ecb.NewECBDecrypter(block)
    decrypter.CryptBlocks(dst, bs)
    fmt.Println(string(dst))

    dst2 := make([]byte, len(bs))
    encrypter := ecb.NewECBEncrypter(block)
    encrypter.CryptBlocks(dst2, dst)
    fmt.Println(base64.StdEncoding.EncodeToString(dst2))
}
