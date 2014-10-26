package pkcs7
import (
)

// ought this to write from src to a
// dst buffer like stdlib stuff seems
// to do?
func Pad(bs []byte, mult int) []byte {
    length := len(bs)

    rem := length % mult
    if rem == 0 {
        return bs
    }

    diff := mult - rem
    res := make([]byte, length + diff)

    copy(res, bs)

    for i := 0; i < diff; i++ {
        res[length + i] = byte(diff)
    }

    return res

}
