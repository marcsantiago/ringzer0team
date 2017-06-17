package main

import (
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"../auth"
)

func binToString(s []byte) string {
	output := make([]byte, len(s)/8)
	for i := 0; i < len(output); i++ {
		val, err := strconv.ParseInt(string(s[i*8:(i+1)*8]), 2, 64)
		if err == nil {
			output[i] = byte(val)
		}
	}
	return string(output)
}

func main() {
	c, err := auth.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.Session.Get("https://ringzer0team.com/challenges/14", nil)
	if err != nil {
		log.Fatal(err)
	}

	html := res.String()
	startChunk := strings.Index(html, "----- BEGIN MESSAGE -----")
	endChunk := strings.Index(html, "----- END MESSAGE -----")
	if startChunk == -1 {
		log.Fatalln("Auth might have failed, can't find PEM body")
	}
	pem := html[startChunk:endChunk]
	r := strings.NewReplacer(
		"----- BEGIN MESSAGE -----<br />", "",
		"<br />", "",
		"\n", "",
		"\r", "",
	)
	// clean up
	pem = strings.TrimSpace(r.Replace(pem))

	// binary to text
	pem = binToString([]byte(pem))

	// hash
	h512 := sha512.New()
	io.WriteString(h512, pem)
	ans := fmt.Sprintf("%x", h512.Sum(nil))
	c.SubmitAnswer("14", ans)
}
