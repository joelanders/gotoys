package xor

import (
    "fmt"
    "encoding/hex"
    "sort"
    "io/ioutil"
    "strings"
    "strconv"
    "encoding/base64"
)

func XorByteSlice(b byte, s []byte) []byte {
    result := make([]byte, len(s))
    for i, _ := range(s) {
        result[i] = s[i] ^ b
    }
    return result
}

func XorSlices(a, b []byte) []byte {
    if len(a) != len(b) {
        panic("unequal lengths")
    }
    result := make([]byte, len(a))
    for i, _ := range(a) {
        result[i] = a[i] ^ b[i]
    }
    return result
}

func XorByteHex(b byte, h string) []byte {
    sl, err := hex.DecodeString(h)
    if err != nil {
        panic(err)
    }

    return XorByteSlice(b, sl)
}

func AllAscii(s []byte) bool {
    for _, c := range(s) {
        if c > 127 {
            return false
        }
    }
    return true
}

func AsciiPercent(s []byte) float64 {
    num := 0
    for _, c := range(s) {
        if (c < 128) && (c > 96) {
            num++
        }
    }
    return 100.0*float64(num)/float64(len(s))
}

func CountLowers(s []byte) int {
    n := 0
    for _, c := range(s) {
        if c > 65 && c < 122 {
            n++
        }
    }
    return n
}

//todo duplication
func CountSpaces(s []byte) int {
    n := 0
    for _, c := range(s) {
        if c == 32 {
            n++
        }
    }
    return n
}

type Candidate struct {
    Key byte
    Plaintext []byte
    Cipherhex string
    AllAscii bool
    NumLowers int
    AsciiPercent float64
    NumSpaces int
}

func TryKey(k byte, h string) *Candidate {
    cand := new(Candidate)
    cand.Key = k
    cand.Plaintext = XorByteHex(k, h)
    cand.Cipherhex = h
    cand.AllAscii = AllAscii(cand.Plaintext)
    cand.NumLowers = CountLowers(cand.Plaintext)
    cand.AsciiPercent = AsciiPercent(cand.Plaintext)
    cand.NumSpaces = CountSpaces(cand.Plaintext)
    return cand
}

//todo: nasty duplication
// cipherhex is wrong name
func TryKeyBytes(k byte, bs []byte) *Candidate {
    cand := new(Candidate)
    cand.Key = k
    cand.Plaintext = XorByteSlice(k, bs)
    cand.Cipherhex = string(bs)
    cand.AllAscii = AllAscii(cand.Plaintext)
    cand.NumLowers = CountLowers(cand.Plaintext)
    cand.AsciiPercent = AsciiPercent(cand.Plaintext)
    cand.NumSpaces = CountSpaces(cand.Plaintext)
    return cand
}

func TryAllKeys(h string) []*Candidate {
    cands := make([]*Candidate, 256)
    for i := 0; i <= 255; i++ {
        cands[i] = TryKey(byte(i), h)
    }
    return cands
}

//todo: nasty duplication
func TryAllKeysBytes(bs []byte) []*Candidate {
    cands := make([]*Candidate, 256)
    for i := 0; i <= 255; i++ {
        cands[i] = TryKeyBytes(byte(i), bs)
    }
    return cands
}

type ByAscii []*Candidate

func (cs ByAscii) Len() int {
    return len(cs)
}

func (cs ByAscii) Swap(i, j int) {
    cs[i], cs[j] = cs[j], cs[i]
}

func (cs ByAscii) Less(i, j int) bool {
//    return cs[i].AsciiPercent < cs[j].AsciiPercent
    return cs[i].NumSpaces < cs[j].NumSpaces
}


//todo remove duplication here somehow
type ByLowers []*Candidate

func (cs ByLowers) Len() int {
    return len(cs)
}

func (cs ByLowers) Swap(i, j int) {
    cs[i], cs[j] = cs[j], cs[i]
}

func (cs ByLowers) Less(i, j int) bool {
    return cs[i].NumLowers < cs[j].NumLowers
}

func AsciiCandidates(cs []*Candidate) []*Candidate {
    n := 0
    acs := make([]*Candidate, len(cs))

    for _, c := range(cs) {
        if c.AllAscii {
            acs[n] = c
            n ++
        }
    }
    return acs[:n]
}

func PrintCandidate(c *Candidate) {
    fmt.Printf("%d: %f: %s\n\n", c.Key, c.AsciiPercent, strconv.Quote(string(c.Plaintext)))
}

func PrintCandidates(cs []*Candidate) {
    sort.Sort(ByAscii(cs))
    for _, c := range(cs) {
        PrintCandidate(c)
    }
}

func LoadHexes() []string {
    all, err := ioutil.ReadFile("4.txt")
    if err != nil {
        panic(err)
    }

    hexes := strings.Split(string(all), "\n")

    return hexes
}

func AllCandidates(hs []string) []*Candidate {
    var cands []*Candidate
    for _, h := range(hs) {
        cands = append(cands, TryAllKeys(h)...)
    }
    
    fmt.Println(len(cands))

    return cands
}

func TryKeyOnHexes(k byte, hs []string) {
    fmt.Println(len(hs))
    for _, h := range(hs) {
        s := XorByteHex(k, h)
        if AllAscii(s) {
            fmt.Println(XorByteHex(k, h))
        }
    }
}

func RepXor(key, text []byte) []byte {
    cipherText := make([]byte, len(text))
    for i, _ := range(text) {
        cipherText[i] = text[i] ^ key[i % len(key)]
    }
    return cipherText
}

func HexXor(key, text string) string {
    cBin := RepXor([]byte(key), []byte(text))
    buf := make([]byte, 512)
    n := hex.Encode(buf, cBin)
    return string(buf[:n])
}

func BytesFromFile(f string) []byte {
    s, err := ioutil.ReadFile(f)
    if err != nil {
        panic(err)
    }
    
    //todo: decodestring is more conv.
    bs, err := base64.StdEncoding.DecodeString(string(s))
    if err != nil {
        panic(err)
    }

    return bs
}

