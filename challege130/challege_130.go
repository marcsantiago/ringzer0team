package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/levigross/grequests"

	"golang.org/x/crypto/ssh"

	"./auth"
)

var (
	userRe     = regexp.MustCompile(`User:\s*([^<]+)`)
	passwordRe = regexp.MustCompile(`Password:\s*([^<]+)`)
	sshRe      = regexp.MustCompile(`<div\s*class="download">.+"ssh:\/\/([^"]+)`)
	gamesWon   = 0
)

// KeyPrint ...
func KeyPrint(dialAddr string, addr net.Addr, key ssh.PublicKey) error {
	fmt.Printf("%s %s %s\n", strings.Split(dialAddr, ":")[0], key.Type(), base64.StdEncoding.EncodeToString(key.Marshal()))
	return nil
}

type algo struct {
	a, b, guess int
}

/* Guessing algo
int a = 1, b = n, guess = average of previous answers;
while(guess is wrong) {
    if(guess lower than answer) {a = guess;}
    else if(guess higher than answer) {b = guess;}
    guess = (a+b)/2;
} //Go back to while
*/

func (al *algo) run(isLower bool) string {
	if isLower {
		al.a = al.guess
	} else {
		al.b = al.guess
	}
	al.guess = (al.a + al.b) / 2
	answer := strconv.Itoa(al.guess)
	return answer + "\n"
}

func newAlgo() *algo {
	al := new(algo)
	al.a = 1
	al.b = 30000 // some high number...
	return al
}

func main() {
	c, err := auth.GetSess()
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.Session.Get("https://ringzer0team.com/challenges/130", nil)
	if err != nil {
		log.Fatal(err)
	}

	var user, password string
	html := res.String()
	m := userRe.FindStringSubmatch(html)
	if len(m) == 2 {
		user = m[1]
	}
	m = passwordRe.FindStringSubmatch(html)
	if len(m) == 2 {
		password = m[1]

	}

	// user = "number"
	// password = "Z7IwIMRC2dc764L"
	if user == "" || password == "" {
		log.Fatalln("Couldn't parse user and password out of html")
	}

	var connStr string
	m = sshRe.FindStringSubmatch(html)
	if len(m) == 2 {
		connStr = m[1]
	}
	// connStr = "ringzer0team.com:12643"

	if !strings.Contains(connStr, ":") {
		log.Fatalln("Missing port")
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: KeyPrint,
	}

	connection, err := ssh.Dial("tcp", connStr, sshConfig)
	if err != nil {

		log.Fatalln(err)
	}

	session, err := connection.NewSession()
	if err != nil {
		log.Fatalln(err)
	}

	// ---------------------------------

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		log.Fatalf("request for pseudo terminal failed: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout for session: %v", err)
	}
	// go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Fatalf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)

	go session.Run("")
	driver := newAlgo()
	first := true
	content := make([]byte, 1024)
	var flag string
	for {
		stdout.Read(content)
		lines := string(content)
		fmt.Println(lines, driver)
		if len(content) == 0 {
			fmt.Println("waiting")
		} else if strings.Contains(lines, "too big") {
			g := driver.run(false)
			stdin.Write([]byte(g))
		} else if strings.Contains(lines, "too small") {
			g := driver.run(true)
			stdin.Write([]byte(g))
		} else if strings.Contains(lines, "You got the right number") {
			gamesWon++
			driver = newAlgo()
			first = true
		}

		if len(content) > 0 && first {
			first = false
			g := driver.run(true)
			stdin.Write([]byte(g))
		}
		f, err := c.GetFlag(lines)
		if err == nil {
			flag = f
			break
		}
	}

	res, err = c.Session.Get("https://ringzer0team.com/challenges/130", nil)
	if err != nil {
		log.Fatal(err)
	}
	csrfToken, err := c.GetCSRF(html)
	if err != nil {
		log.Fatalln(err)
	}

	data := map[string]string{"id": "130", "flag": flag, "check": "false", "csrf": csrfToken}
	res, err = c.Session.Post("https://ringzer0team.com/challenges/130", &grequests.RequestOptions{
		Data: data,
	})
	html = res.String()
	if strings.Contains(html, "Wrong flag try harder!") {
		log.Fatalln("Wrong answer", flag)
	}
	log.Println("Answer seems correct")

	return
}
