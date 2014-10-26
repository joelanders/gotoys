package main
import (
    "fmt"
    "io/ioutil"
    "encoding/base64"
    "math"
    "sort"
    "./xor"
)

func main() {
    fmt.Println(hamByte(2, 2))
    fmt.Println(hamString("this is a test", "wokka wokka!!!"))

    bs := bytesFromFile("6.txt")
    fmt.Println(len(bs))

//    tryKeySizes(2, 40, bs)

//    fmt.Println(transpose(2, []byte{65, 66, 67, 68}))
//    fmt.Println(transpose(3, []byte{65, 66, 67, 68,
//                                    69, 70, 71, 72, 73}))

    indies := transpose(29, bs)
    key := make([]byte, len(indies))

    for i := 0; i < len(indies); i++ {
        cands := xor.TryAllKeysBytes(indies[i])
        top := TopCandidate(cands)
        xor.PrintCandidate(top)
        key[i] = top.Key
    }

    fmt.Println(string(key))
    fmt.Println(string(xor.RepXor(key, bs)))
//    asciiCands := xor.AsciiCandidates(cands)
//    fmt.Println(len(asciiCands))
//    fmt.Println(len(cands))
//    xor.PrintCandidates(asciiCands)
//    for _, cand := range(asciiCands) {
//        fmt.Println(string(cand.Plaintext[:30]))
//    }
}

func TopCandidate(cands []*xor.Candidate) *xor.Candidate {
    sort.Sort(xor.ByAscii(cands))
    return cands[len(cands)-1]
}

func hamByte(a, b byte) int {
    dist := 0
    val := a ^ b
    for val > 0 {
        dist++
        val &= val - 1
    }
    return dist
}

func hamBytes(a, b []byte) int {
    if len(a) != len(b) {
        panic("lengths must equal")
    }
    dist := 0
    for i, _ := range(a) {
        dist += hamByte(a[i], b[i])
    }
    return dist
}

func hamString(a, b string) int {
    return hamBytes([]byte(a), []byte(b))
}

//todo: probe more bytes
func tryKeySize(n int, text []byte) float64 {
    comps := 10
    ham := 0
    for i := 0; i < comps; i++ {
        ham += hamBytes(text[(i+0)*n:(i+1)*n],
                        text[(i+1)*n:(i+2)*n])
    }
    ave := float64(ham) / float64(comps)
    return ave/float64(n)
}

func tryKeySizes(start, end int, text []byte) {
    for i := start; i <= end; i++ {
        fmt.Printf("%d: %f\n", i, tryKeySize(i, text))
    }
}

func bytesFromFile(f string) []byte {
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

func transpose(n int, bs []byte) [][]byte {
    res := make([][]byte, n)
    subLen := int(math.Ceil(float64(len(bs)) / float64(n)))
    for i := range(res) {
        res[i] = make([]byte, subLen)
    }

    for i := 0; i < n; i++ {
        for j := 0; (j*n)+i < len(bs); j++ {
            res[i][j] = bs[(j*n)+i]
        }
    }
    return res
}

