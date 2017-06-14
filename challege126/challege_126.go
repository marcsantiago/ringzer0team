package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"../auth"

	"github.com/levigross/grequests"
)

// personal gist containing a dictionary of words
var gist = "https://api.github.com/gists/a4f87dc7178f1e2c134c82bfda7fbba0"

// Gist response
type Gist struct {
	Comments    int       `json:"comments"`
	CommentsURL string    `json:"comments_url"`
	CommitsURL  string    `json:"commits_url"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	Files       struct {
		DictionaryTxt struct {
			Content   string `json:"content"`
			Filename  string `json:"filename"`
			Language  string `json:"language"`
			RawURL    string `json:"raw_url"`
			Size      int    `json:"size"`
			Truncated bool   `json:"truncated"`
			Type      string `json:"type"`
		} `json:"dictionary2.txt"`
	} `json:"files"`
	Forks      []interface{} `json:"forks"`
	ForksURL   string        `json:"forks_url"`
	GitPullURL string        `json:"git_pull_url"`
	GitPushURL string        `json:"git_push_url"`
	History    []struct {
		ChangeStatus struct {
			Additions int `json:"additions"`
			Deletions int `json:"deletions"`
			Total     int `json:"total"`
		} `json:"change_status"`
		CommittedAt time.Time `json:"committed_at"`
		URL         string    `json:"url"`
		User        struct {
			AvatarURL         string `json:"avatar_url"`
			EventsURL         string `json:"events_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			GravatarID        string `json:"gravatar_id"`
			HTMLURL           string `json:"html_url"`
			ID                int    `json:"id"`
			Login             string `json:"login"`
			OrganizationsURL  string `json:"organizations_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			ReposURL          string `json:"repos_url"`
			SiteAdmin         bool   `json:"site_admin"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			Type              string `json:"type"`
			URL               string `json:"url"`
		} `json:"user"`
		Version string `json:"version"`
	} `json:"history"`
	HTMLURL string `json:"html_url"`
	ID      string `json:"id"`
	Owner   struct {
		AvatarURL         string `json:"avatar_url"`
		EventsURL         string `json:"events_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		GravatarID        string `json:"gravatar_id"`
		HTMLURL           string `json:"html_url"`
		ID                int    `json:"id"`
		Login             string `json:"login"`
		OrganizationsURL  string `json:"organizations_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		ReposURL          string `json:"repos_url"`
		SiteAdmin         bool   `json:"site_admin"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		URL               string `json:"url"`
	} `json:"owner"`
	Public    bool        `json:"public"`
	Truncated bool        `json:"truncated"`
	UpdatedAt time.Time   `json:"updated_at"`
	URL       string      `json:"url"`
	User      interface{} `json:"user"`
}

type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortRunes) Len() int {
	return len(s)
}

func sortString(s string) string {
	r := []rune(s)
	sort.Sort(sortRunes(r))
	return string(r)
}

func correctWord(input, dictionaryWord string) (result string) {
	sInput := sortString(input)
	sDict := sortString(dictionaryWord)
	if sInput == sDict {
		result = dictionaryWord
	}
	return
}

func main() {
	// get dictionary of words
	res, err := grequests.Get(gist, nil)
	if err != nil {
		log.Fatal(err)
	}
	var data Gist
	res.JSON(&data)
	var content string
	if data.Files.DictionaryTxt.Truncated {
		res, err := grequests.Get(data.Files.DictionaryTxt.RawURL, nil)
		if err != nil {
			log.Fatal(err)
		}

		content = res.String()
	}

	var words []string
	if content == "" {
		words = strings.Split(data.Files.DictionaryTxt.Content, "\n")
	} else {
		words = strings.Split(content, "\n")
	}

	// using more mem... give me more speed
	wordSet := make(map[string]bool)
	for i, w := range words {
		// make sure the data is cleanly inputted, though it should be
		clean := strings.TrimSpace(strings.ToLower(w))
		words[i] = clean
		wordSet[clean] = true
	}

	c, err := auth.GetSess()
	if err != nil {
		log.Fatal(err)
	}

	res, err = c.Session.Get("https://ringzer0team.com/challenges/126", nil)
	if err != nil {
		log.Fatal(err)
	}
	html := res.String()
	startChunk := strings.Index(html, "----- BEGIN WORDS -----")
	endChunk := strings.Index(html, "----- END WORDS -----")
	if startChunk == -1 {
		log.Fatalln("can't get word list")
	}
	rawWordsList := html[startChunk:endChunk]

	r := strings.NewReplacer(
		"----- BEGIN WORDS -----<br />", "",
		"<br />", "",
		"\n", "",
		"\r", "",
	)
	// clean up
	rawWordsList = strings.TrimSpace(r.Replace(rawWordsList))
	testWords := strings.Split(rawWordsList, ",")

	var possbleWords []string
	for _, test := range testWords {
		if _, ok := wordSet[test]; ok {
			possbleWords = append(possbleWords, test)
			continue
		}
		for _, w := range words {
			a := correctWord(test, w)
			if a != "" {
				possbleWords = append(possbleWords, a)
				break
			}

		}
	}

	if len(testWords) != len(possbleWords) {
		log.Fatalf("Length of words from site: %d doesn't match the length of words matched: %d", len(testWords), len(possbleWords))
	}

	answer := strings.Join(possbleWords, ",")

	u := fmt.Sprintf("https://ringzer0team.com/challenges/126/%s", answer)
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
	d := map[string]string{"id": "126", "flag": flag, "check": "false", "csrf": csrfToken}
	res, err = c.Session.Post("https://ringzer0team.com/challenges/126", &grequests.RequestOptions{
		Data: d,
	})
	html = res.String()
	if strings.Contains(html, "Wrong flag try harder!") {
		log.Fatalln("Wrong answer")
	}
	log.Println("Answer seems correct")

}
