package main

import (
    "fmt"
    "encoding/hex"
    "sort"
    "io/ioutil"
    "strings"
)

func main() {
//    s := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
//    cands := tryAllKeys(s)
    hexes := loadHexes()
    fmt.Printf("%v", hexes)
    fmt.Println("PRICK")
    fmt.Println(len(hexes))

    cands := allCandidates(hexes)
    fmt.Println("SHIT")
    fmt.Println(len(cands))

    asciiCands := asciiCandidates(cands)
    fmt.Println("FUCK")
    fmt.Println(len(asciiCands))
    printCandidates(asciiCands)

//    tryKeyOnHexes(byte(53), hexes)
}

func xorByteSlice(b byte, s []byte) []byte {
    result := make([]byte, len(s))
    for i, _ := range(s) {
        result[i] = s[i] ^ b
    }
    return result
}

func xorByteHex(b byte, h string) []byte {
    sl, err := hex.DecodeString(h)
    if err != nil {
        panic(err)
    }

    return xorByteSlice(b, sl)
}

func allAscii(s []byte) bool {
    for _, c := range(s) {
        if c > 127 {
            return false
        }
    }
    return true
}

func countLowers(s []byte) int {
    n := 0
    for _, c := range(s) {
        if c > 65 && c < 122 {
            n++
        }
    }
    return n
}

type candidate struct {
    key byte
    plaintext []byte
    cipherhex string
    allAscii bool
    numLowers int
}

func tryKey(k byte, h string) *candidate {
    cand := new(candidate)
    cand.key = k
    cand.plaintext = xorByteHex(k, h)
    cand.cipherhex = h
    cand.allAscii = allAscii(cand.plaintext)
    cand.numLowers = countLowers(cand.plaintext)
    return cand
}

func tryAllKeys(h string) []*candidate {
    cands := make([]*candidate, 256)
    for i := 0; i <= 255; i++ {
        cands[i] = tryKey(byte(i), h)
    }
    return cands
}

type byLowers []*candidate

func (cs byLowers) Len() int {
    return len(cs)
}

func (cs byLowers) Swap(i, j int) {
    cs[i], cs[j] = cs[j], cs[i]
}

func (cs byLowers) Less(i, j int) bool {
    return cs[i].numLowers < cs[j].numLowers
}

func asciiCandidates(cs []*candidate) []*candidate {
    n := 0
    acs := make([]*candidate, len(cs))

    for _, c := range(cs) {
        if c.allAscii {
            acs[n] = c
            n ++
        }
    }
    return acs[:n]
}


func printCandidates(cs []*candidate) {
    sort.Sort(byLowers(cs))
    for _, c := range(cs) {
        fmt.Printf("%d: %d: %s\n", c.key, c.numLowers, c.plaintext)
    }
}

func loadHexes() []string {
    all, err := ioutil.ReadFile("4.txt")
    if err != nil {
        panic(err)
    }

    hexes := strings.Split(string(all), "\n")

    return hexes
}

func allCandidates(hs []string) []*candidate {
    var cands []*candidate
    for _, h := range(hs) {
        cands = append(cands, tryAllKeys(h)...)
    }
    
    fmt.Println(len(cands))

    return cands
}

func tryKeyOnHexes(k byte, hs []string) {
    fmt.Println(len(hs))
    for _, h := range(hs) {
        s := xorByteHex(k, h)
        if allAscii(s) {
            fmt.Println(xorByteHex(k, h))
        }
    }
}
