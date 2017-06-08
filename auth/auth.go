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
func (c Client) GetFlag(html string) (flag string, err error) {
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
func (c Client) GetCSRF(html string) (csrf string, err error) {
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

// GetSess ...
func GetSess() (c Client, err error) {
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
	csrfToken, errTemp := c.GetCSRF(html)
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
