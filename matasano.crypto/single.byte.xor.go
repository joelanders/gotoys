package main

import (
    "fmt"
    "./xor"
)

func main() {
//    s := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
//    cands := xor.TryAllKeys(s)
    hexes := xor.LoadHexes()
//    fmt.Printf("%v", hexes)
    fmt.Println(len(hexes))

    cands := xor.AllCandidates(hexes)
    fmt.Println(len(cands))

    asciiCands := xor.AsciiCandidates(cands)
    fmt.Println(len(asciiCands))
    xor.PrintCandidates(asciiCands)

//    tryKeyOnHexes(byte(53), hexes)
}

