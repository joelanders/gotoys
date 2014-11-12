package pkcs7
import (
    "errors"
)

// ought this to write from src to a
// dst buffer like stdlib stuff seems
// to do?
func Pad(bs []byte, mult int) []byte {
    length := len(bs)

    rem := length % mult

    diff := mult - rem
    res := make([]byte, length + diff)

    copy(res, bs)

    for i := 0; i < diff; i++ {
        res[length + i] = byte(diff)
    }

    return res

}

//todo: should i give mult?
//todo: be more strict
func Unpad(bs []byte) ([]byte, error) {
//    if len(bs) % mult != 0 {
//        return []byte{}, errors.New("not blocksize multiple")
//    }

    padLen := int(bs[len(bs)-1])
    for _, c := range(bs[len(bs) - padLen : len(bs)]) {
        if c != byte(padLen) {
            return []byte{}, errors.New("malformed padding")
        }
    }

    return bs[0:len(bs) - padLen], nil
}
