package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/levigross/grequests"
)

var (
	loginPage = "https://ringzer0team.com/login"
	username  string
	password  string
	parseCSRF = regexp.MustCompile(`var\s*_.+\s*=.+'(.+)';`)
	parseFlag = regexp.MustCompile(`(FLAG-\s*[^<\/]+)`)
)

func init() {
	username = os.Getenv("username")
	password = os.Getenv("password")
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		log.Fatalln("username or password is emtpy in os enviroment")
	}
}

// Client ...
type Client struct {
	Session *grequests.Session
}

// GetFlag ...
func GetFlag(html string) (flag string, err error) {
	m := parseFlag.FindStringSubmatch(html)
	if len(m) == 2 {
		flag = m[1]
	}
	if flag == "" {
		err = fmt.Errorf("flag token is empty")
		return
	}
	return
}

// GetCSRF ...
func GetCSRF(html string) (csrf string, err error) {
	m := parseCSRF.FindStringSubmatch(html)
	if len(m) == 2 {
		csrf = m[1]
	}
	if csrf == "" {
		err = fmt.Errorf("csrf token is empty")
		return
	}
	return
}

// NewSession ...
func NewSession() (c Client, err error) {
	cj, _ := cookiejar.New(nil)

	c.Session = grequests.NewSession(&grequests.RequestOptions{
		UseCookieJar:  true,
		RedirectLimit: 0,
		Host:          "ringzer0team.com",
		CookieJar:     cj,
		Headers: map[string]string{
			"Referer":    "https://ringzer0team.com/login",
			"Origin":     "https://ringzer0team.com",
			"DNT":        "1",
			"Connection": "keep-alive",
			"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		},
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	})

	res, err := c.Session.Get(loginPage, nil)
	if err != nil {
		err = fmt.Errorf("%v", err)
		return
	}

	defer res.Close()
	if err != nil {
		err = fmt.Errorf("couldnt read body %v", err)
		return
	}
	html := string(res.String())
	csrfToken, errTemp := GetCSRF(html)
	if err != nil {
		err = errTemp
		return
	}

	data := map[string]string{"username": username, "password": password, "csrf": csrfToken, "check": "false"}
	res, err = c.Session.Post(loginPage, &grequests.RequestOptions{
		Data:         data,
		UseCookieJar: true,
	})
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code not 200: %d", res.StatusCode)
	}

	base, _ := url.Parse(loginPage)
	if len(c.Session.RequestOptions.CookieJar.Cookies(base)) == 0 {
		err = fmt.Errorf("No cookies Stored")
	}
	if err != nil {
		err = fmt.Errorf("%v", err)
		return
	}
	return
}

// SubmitAnswer ...
func (c Client) SubmitAnswer(challegeNumber, answer string) (err error) {
	page := "https://ringzer0team.com/challenges"
	// get flag page
	flagURL := fmt.Sprintf("%s/%s/%s", page, challegeNumber, answer)
	postURL := fmt.Sprintf("%s/%s", page, challegeNumber)

	res, err := c.Session.Get(flagURL, nil)
	if err != nil {
		return
	}

	// parse flag
	html := res.String()
	flag, err := GetFlag(html)
	if err != nil {
		return
	}
	csrfToken, err := GetCSRF(html)
	if err != nil {
		return
	}
	// post the flag back
	data := map[string]string{"id": challegeNumber, "flag": flag, "check": "false", "csrf": csrfToken}
	res, err = c.Session.Post(postURL, &grequests.RequestOptions{
		Data: data,
	})
	html = res.String()
	if strings.Contains(html, "Wrong flag try harder!") {
		return fmt.Errorf("Answer seems wrong")
	}
	log.Println("Answer seems correct")
	return nil
}
