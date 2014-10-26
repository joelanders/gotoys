package main
import (
    "fmt"
    "./xor"
    "crypto/aes"
    "crypto/cipher"
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
    decrypter := NewCBCDecrypter(block, iv)
    decrypter.CryptBlocks(dst, bs)
    fmt.Println(string(dst))

    dst2 := make([]byte, len(bs))
    encrypter := NewCBCEncrypter(block, iv)
    encrypter.CryptBlocks(dst2, dst)
    fmt.Println(base64.StdEncoding.EncodeToString(dst2))
}

func NewCBCEncrypter(b cipher.Block, iv []byte) CBCEncrypter {
    return CBCEncrypter{b, iv}
}

func NewCBCDecrypter(b cipher.Block, iv []byte) CBCDecrypter {
    return CBCDecrypter{b, iv}
}

type CBCEncrypter struct {
    cipher.Block
    iv []byte
}

type CBCDecrypter struct {
    cipher.Block
    iv []byte
}

//todo dupes
func (crypter CBCEncrypter) CryptBlocks(dst, src []byte) {
    bs := crypter.BlockSize()
    if len(src) % crypter.BlockSize() != 0 {
        panic("bad src length")
    }

    preXor := crypter.iv

    for i := 0; i < len(src); i = i + bs {
        input := xor.XorSlices(preXor, src[i:i+bs])
        crypter.Encrypt(dst[i:i+bs], input)
        preXor = dst[i:i+bs]
    }
}

func (crypter CBCDecrypter) CryptBlocks(dst, src []byte) {
    bs := crypter.BlockSize()
    if len(src) % crypter.BlockSize() != 0 {
        panic("bad src length")
    }

    postXor := crypter.iv

    output := make([]byte, bs)
    for i := 0; i < len(src); i = i + bs {
        crypter.Decrypt(output, src[i:i+bs])
        copy(dst[i:i+bs], xor.XorSlices(postXor, output))
        postXor = src[i:i+bs]
    }
}
