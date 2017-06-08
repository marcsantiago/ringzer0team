package main

import (
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/levigross/grequests"

	"./auth"
)

var baseURL = "https://ringzer0team.com/challenges/13/%s"

func main() {
	c, err := auth.GetSess()
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
	// get flag page
	u := fmt.Sprintf("https://ringzer0team.com/challenges/13/%s", ans)
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
	data := map[string]string{"id": "13", "flag": flag, "check": "false", "csrf": csrfToken}
	res, err = c.Session.Post("https://ringzer0team.com/challenges/13", &grequests.RequestOptions{
		Data: data,
	})
	html = res.String()
	if strings.Contains(html, "Wrong flag try harder!") {
		log.Fatalln("Wrong answer")
	}
	log.Println("Answer seems correct")
}
