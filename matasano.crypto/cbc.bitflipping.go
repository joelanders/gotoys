package main
import (
    "fmt"
    "./cbc"
    "crypto/aes"
    "crypto/cipher"
    "./pkcs7"
    "./xor"
)

const (
    prefix = "comment1=cooking%20MCs;userdata="
    suffix = ";comment2=%20like%20a%20pound%20of%20bacon"
)

var encrypter cipher.BlockMode
var decrypter cipher.BlockMode

func main() {
    InitializeCiphers()

    mid := []byte("pwn")
    plain := CatCookie(mid)
    padCookie := pkcs7.Pad(plain, aes.BlockSize)

    cipherCookie := EncCookie(mid)

    desiredCookie := make([]byte, len(padCookie))
    copy(desiredCookie, padCookie)
    adString := []byte(";;admin=true;pwn")
    copy(desiredCookie[48:], adString)
    fmt.Println(string(desiredCookie))

    xorDiff := xor.XorSlices(padCookie, desiredCookie)
    // shift the xor diff one block to the left
    xorDiff = append(xorDiff, make([]byte, 16)...)
    xorDiff = xorDiff[aes.BlockSize:]
    fmt.Println(xorDiff)

    fakeCookie := xor.XorSlices(xorDiff, cipherCookie)

//     "comment1=cooking %20MCs;userdata= pwn;comment2=%20 like%20a%20pound %20of%20bacon***"
//     "comment1=cooking %20MCs;userdata= pwn;             ;;admin=true;pwn %20of%20bacon***%


    pt, err := DecCookie(fakeCookie)
    if err != nil {
        panic(err)
    }
    fmt.Println(pt)

}

func CatCookie(middle []byte) []byte {
    preAndMid := append([]byte(prefix), middle...)
    plaintext := append(preAndMid, []byte(suffix)...)
    return plaintext
}

func InitializeCiphers() {
    key := []byte("YELLOW SUBMARINE")
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }

    iv := make([]byte, 16)
    encrypter = cbc.NewCBCEncrypter(block, iv)
    decrypter = cbc.NewCBCDecrypter(block, iv)
}

func EncCookie(middle []byte) []byte {
    plaintext := CatCookie(middle)
    padded := pkcs7.Pad(plaintext, aes.BlockSize)
    dst := make([]byte, len(padded))
    encrypter.CryptBlocks(dst, padded)
    return dst
}

func DecCookie(cookie []byte) (string, error) {
    dst := make([]byte, len(cookie))
    decrypter.CryptBlocks(dst, cookie)
    unp, err := pkcs7.Unpad(dst)
    if err != nil {
        return string(unp), err
    }
    return string(unp), nil
}
