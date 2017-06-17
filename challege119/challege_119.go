package main

import (
	"fmt"
	"log"
	"strings"

	"../auth"
)

var (
	nMap = map[string]string{
		"xxx x   xx   xx   x xxx":      "0",
		"xx  x x    x    x  xxxxx":     "1",
		"xxx x   x   xx  x   xxxxx":    "2",
		"xxx x   x  xx x   x xxx":      "3",
		"x   xx    x xxxxx     x    x": "4",
		"xxxxxx     xxxx    xxxxxx":    "5",
	}
)

func main() {
	c, err := auth.NewSession()
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
	c.SubmitAnswer("119", num)
}
