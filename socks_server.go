package main

import (
    "fmt"
    "io"
    "log"
    "net"
    "os"
    "encoding/binary"
    "sync"
    "errors"
)

func HandleCon(locCon net.Conn) {
    // application sends us a SOCKS4 request
    req, err := readConnectRequest(locCon)
    if err != nil {
        log.Println(err)
        return
    }

    // we send them back an "approved" response
    err = sendConnectResp(locCon)
    if err != nil {
        log.Println(err)
        return
    }

    log.Println(req.String())

    // dial the requested destination
    remCon, err := net.Dial("tcp", req.DestAddr())
    if err != nil {
        log.Println(err)
        return
    }
    defer remCon.Close()

    var localReader, remoteReader io.Reader
    // if plaintext, tee to stdout
    if req.port == 80 {
        // localReader is a Reader that acts like locCon,
        // except it also Writes to os.Stdout when we Read from it.
        localReader = io.TeeReader(locCon, os.Stdout)
        remoteReader = io.TeeReader(remCon, os.Stdout)
    } else {
        localReader = locCon
        remoteReader = remCon
     }

    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        defer wg.Done()
        io.Copy(remCon, localReader)
    }()
    go func() {
        defer wg.Done()
        io.Copy(locCon, remoteReader)
    }()
    wg.Wait()
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

func readConnectRequest(c net.Conn) (conReq *ConReq, err error) {
    // read first 9 bytes from the connection
    var b [9]byte
    n, err := c.Read(b[:])
    if n != 9 || err != nil {
        return nil, errors.New("not enough in req")
    }

    // only socks v4 for now
    if b[0] != '\x04' {
        return nil, errors.New("not socks 4")
    }

    req := new(ConReq)

    if b[1] != 1 && b[1] != 2 {
        return nil, errors.New("bad command")
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
    // Listen on TCP port 2000 on loopback interface
    l, err := net.Listen("tcp", "127.0.0.1:2000")
    if err != nil {
        log.Fatal(err)
    }
    defer l.Close()
    for {
        conn, err := l.Accept()
        if err != nil {
            log.Println(err)
        }
        go func() {
            defer conn.Close()
            HandleCon(conn)
        }()
    }
}
