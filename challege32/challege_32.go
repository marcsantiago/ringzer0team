package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"../auth"
)

var parseEquation = regexp.MustCompile(`(\d+)\s*\+\s*0x([0-9-a-z]+)\s*-\s*(\d+)`)

func main() {
	c, err := auth.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.Session.Get("https://ringzer0team.com/challenges/32", nil)
	if err != nil {
		log.Fatal(err)
	}

	html := res.String()
	startChunk := strings.Index(html, "----- BEGIN MESSAGE -----")
	endChunk := strings.Index(html, "----- END MESSAGE -----")
	if startChunk == -1 {
		log.Fatalln("Auth might have failed, can't find PEM body")
	}
	rawEquation := html[startChunk:endChunk]

	var d, h, b int64
	m := parseEquation.FindStringSubmatch(rawEquation)
	if len(m) == 4 {
		tmp, _ := strconv.Atoi(m[1])
		d = int64(tmp)
		h, _ = strconv.ParseInt(m[2], 16, 64)
		b, _ = strconv.ParseInt(m[3], 2, 64)
	}
	answer := d + h - b
	c.SubmitAnswer("32", fmt.Sprintf("%d", answer))
}
