package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"

	"./auth"
)

var (
	userRe     = regexp.MustCompile(`User:\s*([^<]+)`)
	passwordRe = regexp.MustCompile(`Password:\s*([^<]+)`)
	sshRe      = regexp.MustCompile(`<div\s*class="download">.+"ssh:\/\/([^"]+)`)
)

// KeyPrint ...
func KeyPrint(dialAddr string, addr net.Addr, key ssh.PublicKey) error {
	fmt.Printf("%s %s %s\n", strings.Split(dialAddr, ":")[0], key.Type(), base64.StdEncoding.EncodeToString(key.Marshal()))
	return nil
}

/* Guessing algo
int a = 1, b = n, guess = average of previous answers;
while(guess is wrong) {

    if(guess lower than answer) {a = guess;}
    else if(guess higher than answer) {b = guess;}

    guess = (a+b)/2;

} //Go back to while
*/

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

	if user == "" || password == "" {
		log.Fatalln("Couldn't parse user and password out of html")
	}

	var connStr string
	m = sshRe.FindStringSubmatch(html)
	if len(m) == 2 {
		connStr = m[1]
	}

	if !strings.Contains(connStr, ":") {
		log.Fatalln("Missing port")
	}
	// parts := strings.Split(connStr, ":")
	// if len(parts) != 2 {
	// 	log.Fatalf("This %s isn't a connection string\n", connStr)
	// }

	// host := parts[0]
	// port := parts[1]

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
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Fatalf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)

	go session.Start("")
	for {
		stdin.Write([]byte("10000000\n"))
		break
	}
	session.Wait()

}
