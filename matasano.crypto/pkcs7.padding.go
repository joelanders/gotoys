package main
import(
    "fmt"
    "strconv"
    "./pkcs7"
)

func main() {
    strings := []string{"joe", "test", "YELLOW SUBMARINE"}

    fmt.Println("padding to block size of 8")
    for _, s := range(strings) {
        fmt.Println(strconv.Quote(string(pkcs7.Pad([]byte(s), 8))))
    }

    fmt.Println("padding to block size of 20")
    for _, s := range(strings) {
        fmt.Println(strconv.Quote(string(pkcs7.Pad([]byte(s), 20))))
    }

}

