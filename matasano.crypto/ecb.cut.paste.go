package main

import (
    "math/rand"
    "bytes"
    "time"
    "strings"
    "strconv"
    "fmt"
    "errors"
    "crypto/aes"
    "crypto/cipher"
    "./ecb"
    "./pkcs7"
)

const (
)

var encrypter cipher.BlockMode
var decrypter cipher.BlockMode

func main() {
    key := []byte("yellow submarine")

    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }

    encrypter = ecb.NewECBEncrypter(block)
    decrypter = ecb.NewECBDecrypter(block)

    rand.Seed( time.Now().UTC().UnixNano())

    // we'll request the encrypted cookies of a real admin and a real user,
    // two at a time. cut and paste the last block of the admin into the
    // last block of the user, then submit this fake cookie and see if it
    // decrypts to a valid user (with role=admin).
    for {
        //simulate our stolen encrypted cookie
        adminCookie := randCookie("admin")
        cipherAdmin := encryptString(adminCookie)

        userCookie := randCookie("user")
        cipherUser := encryptString(userCookie)

        // just for laziness later
        if len(cipherAdmin) != 48 || len(cipherUser) != 48 {
            continue
        }

        fmt.Println("######################")

        fmt.Println(userFromEncryptedCookie(string(cipherAdmin)))
        fmt.Println(userFromEncryptedCookie(string(cipherUser)))
        fmt.Println()

        fakeCookie := pasteLastBlock(cipherUser, cipherAdmin)
        fakeUser, err := userFromEncryptedCookie(string(fakeCookie))
        if err == nil {
            fmt.Println("good fake user!")
            fmt.Println("real admin cookie:")
            fmt.Println(cipherAdmin)
            fmt.Println("real user cookie:")
            fmt.Println(cipherUser)
            fmt.Println("fake user cookie:")
            fmt.Println(fakeCookie)
            fmt.Println("fake user:")
            fmt.Println(fakeUser)
            break
        }
    }

}

func randLetter() string {
    asciiLowers := []byte("abcdefghijklmnopqrstuvwxyz")
    return string(asciiLowers[rand.Int()%26])
}

func randString(length int) string {
    var b bytes.Buffer
    for i := 0; i < length; i++ {
        b.WriteString(randLetter())
    }
    return b.String()
}

func randEmail() string {
    len1 := 1 + rand.Int() % 10
    len2 := 1 + rand.Int() % 10

    var email bytes.Buffer
    email.WriteString(randString(len1))
    email.WriteString("@")
    email.WriteString(randString(len2))
    email.WriteString(".com")

    return email.String()
}

type User struct {
    Email string
    Uid int
    Role string
}

//func (u *User) String() string {
//    var str bytes.Buffer
//    str.WriteString("<")
//    str.WriteString(u.Email)
//    str.WriteString(">")
//    str.WriteString(" <")
//    str.WriteString(strconv.Itoa(u.Uid))
//    str.WriteString(">")
//    str.WriteString(" <")
//    str.WriteString(u.Role)
//    str.WriteString(">")
//    return str.String()
//
//}

func userFromCookie(cookie string) (*User, error) {
    kvs := strings.Split(cookie, "&")
    u := new(User)

    if len(kvs) != 3 {
        return u, errors.New("bad req")
    }

    emailPair := strings.Split(kvs[0], "=")
    if len(emailPair) != 2 || emailPair[0] != "email" {
        return u, errors.New("bad email")
    }

    uidPair := strings.Split(kvs[1], "=")
    if len(uidPair) != 2 || uidPair[0] != "uid" {
        return u, errors.New("bad uid")
    }

    rolePair := strings.Split(kvs[2], "=")
    if len(rolePair) != 2 || rolePair[0] != "role" {
        return u, errors.New("bad role")
    }

    u.Email = emailPair[1]
    uid, err := strconv.Atoi(uidPair[1])
    if err != nil {
        return u, errors.New("bad uid parsed")
    }
    u.Uid = uid
    u.Role = rolePair[1]
    return u, nil
}

func randCookie(role string) string {
    if role != "user" && role != "admin" {
        panic("bad role")
    }
    var cookie bytes.Buffer
    cookie.WriteString("email=")
    cookie.WriteString(randEmail())
    cookie.WriteString("&uid=")
    cookie.WriteString(strconv.Itoa(rand.Int()%10000))
    cookie.WriteString("&role=")
    cookie.WriteString(role)
    return cookie.String()
}

func encryptString(cookie string) []byte {
    input := pkcs7.Pad([]byte(cookie), aes.BlockSize)
    output := make([]byte, len(input))
    encrypter.CryptBlocks(output, input)
    return output
}

func decryptString(cipher string) ([]byte, error) {
    output := make([]byte, len(cipher))
    decrypter.CryptBlocks(output, []byte(cipher))
    cookie, err := pkcs7.Unpad(output, aes.BlockSize)
    if err != nil {
        return []byte{}, errors.New("bad padding")
    }
    return cookie, nil
}

func userFromEncryptedCookie(cipher string) (*User, error) {
    decrypted, err := decryptString(cipher)
    if err != nil {
        return nil, errors.New("bad req")
    }

    user, err := userFromCookie(string(decrypted))
    if err != nil {
        return nil, errors.New("bad req")
    }

    return user, nil
}

//return a new copy, not overwriting
func pasteLastBlock(dst, src []byte) []byte {
    result := make([]byte, len(dst))
    copy(result, dst)
    copy(result[32:48], src[32:48])
    return result
}
