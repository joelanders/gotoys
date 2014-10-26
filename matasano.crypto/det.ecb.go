package main
import (
    "fmt"
    "io/ioutil"
    "strings"
    "encoding/hex"
)

func main() {
    hexes := ReadHexesFromFile("8.txt")
    fmt.Println(len(hexes))

    bytes := HexesToBytes(hexes)
    fmt.Println(len(bytes))

    for _, cipher := range(bytes) {
        if CheckDupeBlock(cipher, 16) {
            fmt.Println(hex.EncodeToString(cipher))
        }
    }
}

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
