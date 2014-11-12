package main
import (
    "fmt"
    "encoding/base64"
    "encoding/hex"
    "crypto/cipher"
    "crypto/rand"
    "crypto/aes"
    "math"
    "errors"
    "bytes"
    "io"
    "./ecb"
    "./pkcs7"
)

var _ = rand.Reader
var _ = io.ReadFull

const b64unknown1 = "Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK"
const b64unknown2 = "dGhlcXVpY2ticm93bmZveGp1bXBzb3ZlcnRoZWxhenlkb2cuamFja2Rhd3Nsb3ZlbXlzcGhpbnhvZnF1YXJ0ego="
var unknown1 []byte //only b64-decoding it once
var unknown2 []byte //only b64-decoding it once
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

    unknown1, err = base64.StdEncoding.DecodeString(b64unknown1)
    if err != nil {
        panic(err)
    }

    unknown2, err = base64.StdEncoding.DecodeString(b64unknown2)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(CrackOracle(EncWithPrefix)))
    fmt.Println(string(CrackOracle(EncWithUnknownPrefix)))

}

func CrackOracle(oracle Oracle) []byte {
    knownBytes := []byte{}
    length := len(unknown1)
    deadPrefLen := UnknownPrefixLength(oracle) // zero for oracle #1
    for len(knownBytes) < length {
        nextByte, err := CrackNextByte(deadPrefLen, knownBytes, oracle)
        if err != nil {
            break
        }
        knownBytes = append(knownBytes, nextByte)
    }
    return knownBytes
}

type Oracle func([]byte) []byte

func EncWithPrefix(prefix []byte) []byte {
    input := pkcs7.Pad(append(prefix, unknown1...), aes.BlockSize)
    output := make([]byte, len(input))
    encrypter.CryptBlocks(output, input)
    return output
}

func EncWithUnknownPrefix(mid []byte) []byte {
//    key := make([]byte, 16)
//    io.ReadFull(rand.Reader, key)
    preAndMid := append([]byte("bladshfowafeoislkdf323"), mid...)
    input := pkcs7.Pad(append(preAndMid, unknown2...), aes.BlockSize)
    output := make([]byte, len(input))
    encrypter.CryptBlocks(output, input)
    return output
}

func CrackNextByte(deadPrefLen int, knownBytes []byte, oracle Oracle) (byte, error) {
    // todo: bleh come back to these
    deadBlocks := int(math.Ceil(float64(deadPrefLen)/blocksize))
    deadPad := (blocksize - (deadPrefLen % blocksize)) % blocksize
    blkIndex := deadBlocks + len(knownBytes) / blocksize

    // we make padding such that there is only one unknown byte
    // in the block we're looking at
    padLen := blocksize - (len(knownBytes)%blocksize) - 1
    padding := bytes.Repeat([]byte("a"), deadPad + padLen)

    actual := oracle(padding)
    
    if blkIndex * blocksize > len(actual) {
        return byte(0), errors.New("eof")
    }

    padPlusKnown := append(padding, knownBytes...)

    byteThatGives := make(map[string]byte)

    for i := 0; i < 256; i++ {
        knownPlusGuess := append(padPlusKnown, byte(i))
        res := oracle(knownPlusGuess)
        byteThatGives[HexBlock(blkIndex, res)] = byte(i)
    }

    good, ok := byteThatGives[HexBlock(blkIndex, actual)]
    if !ok {
        return byte(0), errors.New("probably eof")
    }

    return good, nil
    
}

func HexBlock(i int, bs []byte) string {
    return hex.EncodeToString(bs[i*blocksize:(i*blocksize)+blocksize])
}

func UnknownPrefixLength(oracle Oracle) int {
    initDeadBlocks := NumDeadBlocks(oracle([]byte{}), oracle([]byte{1}))
    deadBlocks := 0
    num := 0
    lastCipherText := oracle([]byte(""))
    for {
        padding := bytes.Repeat([]byte("a"), num + 1)
        nextCipherText := oracle(padding)
        deadBlocks = NumDeadBlocks(lastCipherText, nextCipherText)
        if deadBlocks > initDeadBlocks {
            break
        }
        lastCipherText = nextCipherText
        num++
    }
    return 16*initDeadBlocks + (16-num)
}

func NumDeadBlocks(a, b []byte) int {
    i := 0
    for i = 0; i < len(a) && i < len(b); i++ {
        if a[i] != b[i] {
            break
        }
    }
    return i / 16 //integer division
}
