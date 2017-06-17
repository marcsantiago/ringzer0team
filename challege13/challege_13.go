package main

import (
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"strings"

	"../auth"
)

func main() {
	c, err := auth.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.Session.Get("https://ringzer0team.com/challenges/13", nil)
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
	// hash
	h512 := sha512.New()
	io.WriteString(h512, pem)
	ans := fmt.Sprintf("%x", h512.Sum(nil))
	c.SubmitAnswer("13", ans)
}
