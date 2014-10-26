package main
import (
    "fmt"
    "crypto/aes"
    "crypto/cipher"
    //"encoding/hex"
    "./xor" //todo move stuff
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
//    encrypter := NewCBCEncrypter(block)
//    encrypter.CryptBlocks(dst1, src1)
//    fmt.Println(hex.EncodeToString(dst1))

    dst := make([]byte, len(bs))
    decrypter := NewCBCDecrypter(block)
    decrypter.CryptBlocks(dst, bs)
    fmt.Println(string(dst))

}

func NewCBCEncrypter(b cipher.Block) CBCEncrypter {
    return CBCEncrypter{b}
}

func NewCBCDecrypter(b cipher.Block) CBCDecrypter {
    return CBCDecrypter{b}
}

type CBCEncrypter struct {
    cipher.Block
}

type CBCDecrypter struct {
    cipher.Block
}

//todo dupes
func (crypter CBCEncrypter) CryptBlocks(dst, src []byte) {
    bs := crypter.BlockSize()
    if len(src) % crypter.BlockSize() != 0 {
        panic("bad src length")
    }

    for i := 0; i < len(src); i = i + bs {
        crypter.Encrypt(dst[i:i+bs], src[i:i+bs])
    }
}

func (crypter CBCDecrypter) CryptBlocks(dst, src []byte) {
    bs := crypter.BlockSize()
    if len(src) % crypter.BlockSize() != 0 {
        panic("bad src length")
    }

    for i := 0; i < len(src); i = i + bs {
        crypter.Decrypt(dst[i:i+bs], src[i:i+bs])
    }
}
