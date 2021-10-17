package main

import (
    "bytes"
    "os"
    "io"
    "fmt"
    "time"
)

func CopyOut(r io.Reader) {
    pb, ok := r.(BufPeeker)
    if ok {
        pb.BufPeek()
    }
    io.Copy(os.Stdout, r)
}

type InstrReader struct {
    r io.Reader
    b []byte
    f0 func()
    f1 func()
    t time.Time
    c chan byte
    d chan error
    seen bool
}

func NewInstrReader(r io.Reader) *InstrReader {
    ir := new(InstrReader)
    ir.r = r
    ir.b = make([]byte, 1)
    // ir.c = make(chan byte)
    // ir.d = make(chan error)
    return ir
}

func (ir InstrReader) Read(b []byte) (n int, err error) {
    // defer func() {
    //     if err != nil {
    //         fmt.Printf("%v\n", time.Now().Sub(ir.t))
    //     }
    // }()

    // if !ir.seen {
    //     ir.t = time.Now()
    //     ir.seen = true
    // }
    ir.c = make(chan byte)
    ir.d = make(chan error)

    go func() {
        defer close(ir.c)
        defer close(ir.d)

        n, err = ir.r.Read(ir.b[0:1])
        if n > 0 {
            ir.c <- ir.b[0]
        }
        if err != nil {
            ir.d <- io.EOF
        }
    }()

    select {
    case bc := <-ir.c:
        b[0] = bc
        return 1, nil
    case <-time.After(800 * time.Nanosecond):
        b[0] = byte('x')
        return 1, nil
    case err := <-ir.d:
        return 0, err
    }
}

type BufPeeker interface {
    BufPeek()
}

func (ir InstrReader) BufPeek() {
    fmt.Printf("buf: %v\n", ir.b)
}

func main() {
    slice := []byte("the quick brown fox jumps over the lazy dog\n")

    r := NewInstrReader(
        bytes.NewReader(slice),
    )

    CopyOut(r)
    panic("show stacks")
}
