package main
import (
    "fmt"
    "crypto/cipher"
    "crypto/aes"
    "errors"
    "./pkcs7"
    "./cbc"
)

var encrypter cipher.BlockMode
var decrypter cipher.BlockMode

func main() {
    InitializeCiphers()
    pt := []byte("11223344 1223344 abcdefg ABCDEFG hijklmn OPQRSTU")
    bigKnown := []byte{}
    known := []byte{}
    ct := PadEnc(pt)

    for {
        if len(ct) <= 16 {
            break
        }
        next, err := FindNext(ct, known)
        if err != nil {
            panic(err)
        }
        known = append([]byte{next}, known...)
        if len(known) == aes.BlockSize {
            bigKnown = append(known, bigKnown...)
            ct = ct[:len(ct)-aes.BlockSize]
            known = []byte{}
        }
    }
    fmt.Println(bigKnown)
    fmt.Println(string(bigKnown))
}

func FindNext(cipher, plain []byte) (byte, error) {
    if len(cipher) - len(plain) <= aes.BlockSize {
        return byte(0), errors.New("can't decode first block")
    }
    candidates := []byte{}
    for i:=0; i<255; i++ {
        fiddled, err := XorNextAndPad(cipher, plain, byte(i))
        if err != nil {
            panic(err)
        }
        _, err = DecUnpad(fiddled)
        if err != nil {
            continue
        }
        candidates = append(candidates, byte(i))
    }
    // todo: this looks like shit
    if len(candidates) != 1 {
        if len(plain) != 0 {
            panic("think more")
        }
        if len(candidates) != 2 {
            panic("figure it out")
        }
        for _, c := range(candidates) {
            if c != byte(0) {
                return c ^ byte(len(plain) + 1), nil
            }
        }
        panic("shouldn't be here")
    }
    return candidates[0] ^ byte(len(plain) + 1), nil
}

func XorNextAndPad(cipher0, plain []byte, b byte) ([]byte, error) {
    blockFromEnd := len(plain) / aes.BlockSize
    cipher := make([]byte, len(cipher0) - 16*blockFromEnd)
    copy(cipher, cipher0)
    plain = plain[:len(plain) - 16*blockFromEnd]

    length := len(cipher)
    if length % aes.BlockSize != 0 {
        return cipher, errors.New("bad cipher length")
    }
    position := length - 1 - len(plain) - aes.BlockSize
    if position < 0 {
        return cipher, errors.New("can't change first block")
    }

    // if we know 1 byte of plaintext, next val to brute-force is \x02
    padLen := len(plain)
    for i:=0; i<padLen; i++ {
        cipher[position+1+i] = cipher[position+1+i] ^ plain[i] ^ byte(padLen+1)
    }

    cipher[position] = cipher[position] ^ b
    return cipher, nil

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

func PadEnc(plain []byte) []byte {
    padded := pkcs7.Pad(plain, aes.BlockSize)
    dst := make([]byte, len(padded))
    encrypter.CryptBlocks(dst, padded)
    return dst
}

func DecUnpad(cookie []byte) (string, error) {
    dst := make([]byte, len(cookie))
    decrypter.CryptBlocks(dst, cookie)
    unp, err := pkcs7.Unpad(dst)
    if err != nil {
        return string(unp), err
    }
    return string(unp), nil
}
