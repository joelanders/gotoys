package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"strings"
	"sync"
	"net/http"
)

func fibRet(n int) int {
	if n < 2 {
		return n
	}

	return fibRet(n-1) + fibRet(n-2)
}

func fibChan(n int, a chan int) {
	if n < 2 {
		a <- n
		return
	}

	b := make(chan int)
	go fibChan(n-1, b)

	c := make(chan int)
	go fibChan(n-2, c)

	a <- (<-b) + (<-c)
}

func fibTcp(ctx context.Context, n int, a string) {
	// ctx, task := trace.NewTask(ctx, fmt.Sprintf("fibTcp(%d, %s)", n, a))
	// defer task.End()
	trace.WithRegion(ctx, fmt.Sprintf("fibTcp(%d, %s)", n, a), func() {
		var err error = nil
		defer func() {
			log.Printf("err n: %d, a: %s, err: %v\n", n, a, err)
		}()

		c, err := net.Dial("tcp", a)
		if err != nil {
			return
		}
		defer c.Close()

		if n < 2 {
			c.Write([]byte(fmt.Sprintf("%d\n", n)))
			return
		}

		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return
		}
		defer l.Close()

		go fibTcp(ctx, n-1, l.Addr().String())
		go fibTcp(ctx, n-2, l.Addr().String())

		as := [2]int{}

		var wg sync.WaitGroup
		wg.Add(2)
		for i := 0; i < 2; i++ {
			go func(i int) {
				defer wg.Done()
				ac, err := l.Accept()
				if err != nil {
					return
				}
				a, err := bufio.NewReader(ac).ReadString('\n')
				if err != nil {
					return
				}
				a = strings.TrimSpace(a)
				ai, err := strconv.Atoi(a)
				if err != nil {
					return
				}
				as[i] = ai
			}(i)
		}
		wg.Wait()

		c.Write([]byte(fmt.Sprintf("%d\n", as[0]+as[1])))
	})
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var exTrace = flag.String("trace", "", "write trace to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *exTrace != "" {
		f, err := os.Create(*exTrace)
		if err != nil {
			log.Fatal(err)
		}

		_ = trace.Start(f)
		defer trace.Stop()
	}

	N := 10

	// fmt.Println("fibRet: ", fibRet(N))

	// ac := make(chan int)
	// go fibChan(N, ac)

	// a := <-ac
	// fmt.Println("fibChan: ", a)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	ctx := context.Background()
	go fibTcp(ctx, N, l.Addr().String())

	c, err := l.Accept()
	if err != nil {
		panic(err)
	}
	as, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		panic(err)
	}
	as = strings.TrimSpace(as)
	ai, err := strconv.Atoi(as)
	if err != nil {
		panic(err)
	}

	fmt.Println("fibTcp: ", ai)

    fmt.Println(http.ListenAndServe("localhost:8080", nil))
}
