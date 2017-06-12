package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/levigross/grequests"

	"../auth"
)

func main() {
	hashes := make(map[string]int)
	for i := 0; i < 1000000; i++ {
		hash := sha1.Sum([]byte(strconv.Itoa(i)))
		h := hex.EncodeToString(hash[:])
		hashes[h] = i
	}

	c, err := auth.GetSess()
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

	u := fmt.Sprintf("https://ringzer0team.com/challenges/56/%s", answer)
	res, err = c.Session.Get(u, nil)
	if err != nil {
		log.Fatal(err)
	}
	// parse flag
	html = res.String()
	flag, err := c.GetFlag(html)
	if err != nil {
		log.Fatalln("Couldn't find flag in html")
	}

	csrfToken, err := c.GetCSRF(html)
	if err != nil {
		log.Fatalln(err)
	}

	// post the flag back
	data := map[string]string{"id": "56", "flag": flag, "check": "false", "csrf": csrfToken}
	res, err = c.Session.Post("https://ringzer0team.com/challenges/56", &grequests.RequestOptions{
		Data: data,
	})
	html = res.String()
	if strings.Contains(html, "Wrong flag try harder!") {
		log.Fatalln("Wrong answer")
	}
	log.Println("Answer seems correct")

}
