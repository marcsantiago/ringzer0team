package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func nextPassword(n int, c string) func() string {
	r := []rune(c)
	p := make([]rune, n)
	x := make([]int, len(p))
	return func() string {
		p := p[:len(x)]
		for i, xi := range x {
			p[i] = r[xi]
		}
		for i := len(x) - 1; i >= 0; i-- {
			x[i]++
			if x[i] < len(r) {
				break
			}
			x[i] = 0
			if i <= 0 {
				x = x[0:0]
				break
			}
		}
		return string(p)
	}
}

func main() {
	np := nextPassword(6, "abcdefghijklmnopqrstuvwxyz0123456789")
	var buf bytes.Buffer
	for {
		pwd := np()
		if len(pwd) == 0 {
			f, err := os.OpenFile("passwords.txt", os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}
			if _, err = f.WriteString(buf.String()); err != nil {
				panic(err)
			}
			f.Close()
			buf.Reset()
			break
		}
		h1 := sha1.New()
		io.WriteString(h1, pwd)
		h := fmt.Sprintf("%x", h1.Sum(nil))
		buf.WriteString(fmt.Sprintf("%s:%s\n", h, pwd))
		if buf.Len() > 1024*1024*500 {
			f, err := os.OpenFile("passwords.txt", os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}
			if _, err = f.WriteString(buf.String()); err != nil {
				panic(err)
			}
			f.Close()
			buf.Reset()
		}
	}

}
