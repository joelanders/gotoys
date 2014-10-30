package ecb
import (
    "io/ioutil"
    "strings"
    "encoding/hex"
    "crypto/cipher"
)

func ReadHexesFromFile(filename string) []string {
    bs, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }

    return strings.Split(string(bs), "\n")
}

func HexesToBytes(hexes []string) [][]byte {
    bytes := make([][]byte, len(hexes))
    for i, h := range(hexes) {
        //todo "non-name bytes[i] on left side of :="
        bs, err := hex.DecodeString(h)
        if err != nil {
            panic(err)
        }
        bytes[i] = bs
    }
    return bytes
}

func CheckDupeBlock(bs []byte, size int) bool {
    if len(bs) % size != 0 {
        panic("bad length")
    }

    seen := make(map[string]bool)

    for i := 0; i < len(bs); i += size {
        block := bs[i:i+size]
        if seen[string(block)] {
            return true
        }
        seen[string(block)] = true
    }
    return false
}

func NewECBEncrypter(b cipher.Block) ECBEncrypter {
    return ECBEncrypter{b}
}

func NewECBDecrypter(b cipher.Block) ECBDecrypter {
    return ECBDecrypter{b}
}

type ECBEncrypter struct {
    cipher.Block
}

type ECBDecrypter struct {
    cipher.Block
}

//todo dupes
func (crypter ECBEncrypter) CryptBlocks(dst, src []byte) {
    bs := crypter.BlockSize()
    if len(src) % crypter.BlockSize() != 0 {
        panic("bad src length")
    }

    for i := 0; i < len(src); i = i + bs {
        crypter.Encrypt(dst[i:i+bs], src[i:i+bs])
    }
}

func (crypter ECBDecrypter) CryptBlocks(dst, src []byte) {
    bs := crypter.BlockSize()
    if len(src) % crypter.BlockSize() != 0 {
        panic("bad src length")
    }

    for i := 0; i < len(src); i = i + bs {
        crypter.Decrypt(dst[i:i+bs], src[i:i+bs])
    }
}
