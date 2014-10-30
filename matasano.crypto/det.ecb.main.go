package main
import (
    "fmt"
    "./ecb"
    "encoding/hex"
)

func main() {
    hexes := ecb.ReadHexesFromFile("8.txt")
    fmt.Println(len(hexes))

    bytes := ecb.HexesToBytes(hexes)
    fmt.Println(len(bytes))

    for _, cipher := range(bytes) {
        if ecb.CheckDupeBlock(cipher, 16) {
            fmt.Println(hex.EncodeToString(cipher))
        }
    }
}
