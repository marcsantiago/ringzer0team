package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"../auth"
)

var hashRegex = regexp.MustCompile(`(?s)([a-f0-9]{40}).+([a-f0-9]{64})`)

func main() {
	c, err := auth.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.Session.Get("https://ringzer0team.com/challenges/57", nil)
	if err != nil {
		log.Fatal(err)
	}

	html := res.String()
	startChunk := strings.Index(html, "---- BEGIN HASH -----")
	endChunk := strings.Index(html, "----- END SALT -----")
	if startChunk == -1 {
		log.Fatalln("Auth might have failed, can't find PEM body")
	}
	// hash looks like sha1 probs a number... so lets try a brute force hack

	chunk := html[startChunk:endChunk]
	m := hashRegex.FindStringSubmatch(chunk)
	var hash, salt string
	if len(m) == 3 {
		hash = m[1]
		salt = m[2]
	} else {
		log.Fatalln("seems like regex failed")
	}

	// using some high limit...lets try 10,000
	var answer string
	for i := 0; i <= 10000; i++ {
		h1 := sha1.New()
		digest := fmt.Sprintf("%d%s", i, salt)
		io.WriteString(h1, digest)
		h := fmt.Sprintf("%x", h1.Sum(nil))
		if hash == h {
			answer = fmt.Sprintf("%d", i)
			break
		}
	}
	if answer == "" {
		log.Fatalln("couldnt break the hash...try a higher number")
	}
	c.SubmitAnswer("57", answer)
}
