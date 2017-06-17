package main

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"strconv"
	"strings"

	"../auth"
)

func main() {
	hashes := make(map[string]int)
	for i := 0; i < 1000000; i++ {
		hash := sha1.Sum([]byte(strconv.Itoa(i)))
		h := hex.EncodeToString(hash[:])
		hashes[h] = i
	}

	c, err := auth.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.Session.Get("https://ringzer0team.com/challenges/56", nil)
	if err != nil {
		log.Fatal(err)
	}

	html := res.String()
	startChunk := strings.Index(html, "----- BEGIN HASH -----")
	endChunk := strings.Index(html, "----- END HASH -----")
	if startChunk == -1 {
		log.Fatalln("Auth might have failed, can't find PEM body")
	}
	h := html[startChunk:endChunk]
	r := strings.NewReplacer(
		"----- BEGIN HASH -----<br />", "",
		"<br />", "",
		"\n", "",
		"\r", "",
	)
	// clean up
	h = strings.TrimSpace(r.Replace(h))

	var answer string
	if val, ok := hashes[h]; ok {
		answer = strconv.Itoa(val)
	} else {
		log.Fatalln("not in hash map")
	}
	c.SubmitAnswer("56", answer)
}
