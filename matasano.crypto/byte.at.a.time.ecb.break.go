package main
import (
    "fmt"
    "encoding/base64"
    "encoding/hex"
    "crypto/cipher"
    "crypto/rand"
    "crypto/aes"
    "errors"
    "bytes"
    "io"
    "./ecb"
    "./pkcs7"
)

var _ = rand.Reader
var _ = io.ReadFull

const b64unknown = "Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK"
//const b64unknown = "dGhlcXVpY2ticm93bmZveGp1bXBzb3ZlcnRoZWxhenlkb2cuamFja2Rhd3Nsb3ZlbXlzcGhpbnhvZnF1YXJ0ego="
var unknown []byte //only b64-decoding it once
//var blocksize int //which we will deduce
const blocksize = 16
var encrypter cipher.BlockMode

func main() {
//    key := make([]byte, 16)
//    io.ReadFull(rand.Reader, key)
    key := []byte("yellow submarine")

    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }

    encrypter = ecb.NewECBEncrypter(block)

    unknown, err = base64.StdEncoding.DecodeString(b64unknown)
    if err != nil {
        panic(err)
    }

    length := len(unknown)

    knownBytes := []byte{}
    for len(knownBytes) < length {
        nextByte, err := CrackNextByte(knownBytes)
        if err != nil {
            break
        }
        knownBytes = append(knownBytes, nextByte)
    }
    fmt.Println(string(knownBytes))
}

func EncWithPrefix(prefix []byte) []byte {
    input := pkcs7.Pad(append(prefix, unknown...), aes.BlockSize)
    output := make([]byte, len(input))
    encrypter.CryptBlocks(output, input)
    return output
}

func CrackNextByte(knownBytes []byte) (byte, error) {
    blkIndex := len(knownBytes) / blocksize

    // we make padding such that there is only one unknown byte
    // in the block we're looking at
    padLen := blocksize - (len(knownBytes)%blocksize) - 1
    padding := bytes.Repeat([]byte("a"), padLen)

    actual := EncWithPrefix(padding)
    
    if blkIndex * blocksize > len(actual) {
        return byte(0), errors.New("eof")
    }

    padPlusKnown := append(padding, knownBytes...)

    byteThatGives := make(map[string]byte)

    for i := 0; i < 256; i++ {
        knownPlusGuess := append(padPlusKnown, byte(i))
        res := EncWithPrefix(knownPlusGuess)
        byteThatGives[HexBlock(blkIndex, res)] = byte(i)
    }

    good, ok := byteThatGives[HexBlock(blkIndex, actual)]
    if !ok {
        panic("not found")
    }

    return good, nil
    
}

func HexBlock(i int, bs []byte) string {
    return hex.EncodeToString(bs[i*blocksize:(i*blocksize)+blocksize])
}
