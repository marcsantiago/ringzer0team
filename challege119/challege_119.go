package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/levigross/grequests"

	"./auth"
)

var (
	nMap = map[string]string{
		"xxx x   xx   xx   x xxx":      "0",
		"xx  x x    x    x  xxxxx":     "1",
		"xxx x   x   xx  x   xxxxx":    "2",
		"xxx x   x  xx x   x xxx":      "3",
		"xxxxxx     xxxx    xxxxxx":    "5",
		"x   xx    x xxxxx     x    x": "9",
	}
)

func main() {
	c, err := auth.GetSess()
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.Session.Get("https://ringzer0team.com/challenges/119", nil)
	if err != nil {
		log.Fatal(err)
	}

	html := res.String()
	startChunk := strings.Index(html, "----- BEGIN MESSAGE -----")
	endChunk := strings.Index(html, "----- END MESSAGE -----")
	if startChunk == -1 {
		log.Fatalln("Auth might have failed, can't find PEM body")
	}
	code := html[startChunk:endChunk]
	r := strings.NewReplacer(
		"----- BEGIN MESSAGE -----", "",
		"\n", "",
		"<br />", "\n",
		"&nbsp", " ",
		";", "",
		"<CR>", "",
		"\t\t", "",
		"\r\n", "",
	)
	code = r.Replace(code)

	var rawN, num string
	var counter int

	for _, line := range strings.Split(code, "\n") {
		if strings.Contains(line, "x") {
			counter++
			rawN += line + "\n"
		}
		if counter == 5 {
			// remove the trim stuff later using it because i'm printing it with newline for testing
			fmt.Println(rawN)
			key := strings.TrimSpace(strings.Replace(rawN, "\n", "", -1))
			if val, ok := nMap[key]; ok {
				num += val
			} else {
				log.Fatalf("Missing this number in my map %s\n\n", rawN)
			}
			counter = 0
			rawN = ""
		}
	}

	u := fmt.Sprintf("https://ringzer0team.com/challenges/119/%s", num)
	fmt.Println(u)
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
	d := map[string]string{"id": "119", "flag": flag, "check": "false", "csrf": csrfToken}
	res, err = c.Session.Post("https://ringzer0team.com/challenges/119", &grequests.RequestOptions{
		Data: d,
	})
	html = res.String()
	if strings.Contains(html, "Wrong flag try harder!") {
		log.Fatalln("Wrong answer")
	}
	log.Println("Answer seems correct")

}
