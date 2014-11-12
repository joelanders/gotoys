package main
import(
    "fmt"
    "strconv"
    "./pkcs7"
)

//todo lazy dupes
func main() {
    strings := []string{"joe", "test", "YELLOW SUBMARINE"}

    fmt.Println("(un)padding to block size of 8")
    for _, s := range(strings) {
        padded := pkcs7.Pad([]byte(s), 8)
        fmt.Println("padded: ", strconv.Quote(string(padded)))
        unpadded, err := pkcs7.Unpad([]byte(padded))
        if err != nil {
            panic("unpadding 8 failed")
        }
        fmt.Println("unpadded: ", strconv.Quote(string(unpadded)))
    }

    fmt.Println()
    fmt.Println("(un)padding to block size of 20")
    for _, s := range(strings) {
        padded := pkcs7.Pad([]byte(s), 20)
        fmt.Println("padded: ", strconv.Quote(string(padded)))
        unpadded, err := pkcs7.Unpad([]byte(padded))
        if err != nil {
            panic("unpadding 20 failed")
        }
        fmt.Println("unpadded: ", strconv.Quote(string(unpadded)))
    }

    badPad := []byte{9, 9, 9, 9, 5, 5, 5, 5}
    unp, err := pkcs7.Unpad(badPad)
    if err != nil {
        fmt.Println("good, it failed")
    } else {
        fmt.Println("bad")
        fmt.Println(unp)
    }

}

