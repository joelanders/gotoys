package cbc
import (
    "../xor"
    "crypto/cipher"
)

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
