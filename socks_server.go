package main

import (
    "fmt"
    "io"
    "bufio"
    "log"
    "net"
    "os"
    "encoding/binary"
)

func HandleCon(locCon net.Conn) {
    // application sends us a SOCKS4 request
    req, err := readConnectRequest(locCon)
    if err != nil {
        log.Fatal(err)
    }

    // we send them back an "approved" response
    err = sendConnectResp(locCon)
    if err != nil {
        log.Fatal(err)
    }

    log.Println(req.String())

    // dial the requested destination
    remCon, err := net.Dial("tcp", req.DestAddr())
    if err != nil {
        log.Fatal(err)
    }

    var localReader, remoteReader io.Reader
    // if plaintext, tee to stdout
    if req.port == 80 {
        // localReader is a Reader that acts like locCon,
        // except it also Writes to os.Stdout when we Read from it.
        localReader = io.TeeReader(locCon, os.Stdout)
        remoteReader = io.TeeReader(remCon, os.Stdout)
    } else {
        // if encrypted, don't tee to stdout
        localReader = locCon
        remoteReader = remCon
    }

    go io.Copy(remCon, localReader)
    go io.Copy(locCon, remoteReader)
}

type command int

const (
    CONNECT command = iota
    BIND command = iota
)

type ConReq struct {
    cmd command
    ip net.IP
    port uint16
    user string
    localConn net.Conn
}

func (req *ConReq) DestAddr() string {
    return fmt.Sprintf("%s:%d", req.ip.String(), req.port)
}

func (req *ConReq) String() string {
    clientAddr := req.localConn.RemoteAddr()
    return fmt.Sprintf("%s -> %s", clientAddr, req.DestAddr())
}

// todo: i want c to be a net.Conn and an io.Reader, so new type?
func readConnectRequest(c net.Conn) (conReq *ConReq, err error) {
    r := bufio.NewReader(c)
    
    // read first 8 bytes from the connection
    var b [8]byte
    _, err = io.ReadFull(r, b[:])
    if err != nil {
        log.Fatal(err)
    }

    // only socks v4 for now
    if b[0] != '\x04' {
        log.Fatal("not version 4", b[0])
    }

    req := new(ConReq)

    if b[1] != 1 && b[1] != 2 {
        log.Fatal("invalid command", b[1])
    }

    req.cmd = command(b[1])
    req.port = binary.BigEndian.Uint16(b[2:4])
    req.ip = net.IPv4(b[4],b[5],b[6],b[7])
    req.localConn = c

    return req, nil
}

func sendConnectResp(c net.Conn) error {
    // hard-coding success for now
    resp := [8]byte{0, '\x5a', 0, 0, 0, 0, 0, 0}
    c.Write(resp[:])
    return nil
}

func main() {
    // Listen on TCP port 2000 on all interfaces.
    l, err := net.Listen("tcp", ":2000")
    if err != nil {
        log.Fatal(err)
    }
    defer l.Close()
    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatal(err)
        }
        defer conn.Close()
        go HandleCon(conn)
    }
}
