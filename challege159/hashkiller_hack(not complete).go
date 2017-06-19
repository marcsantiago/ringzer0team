package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/levigross/grequests"
	"github.com/nfnt/resize"
	"github.com/otiai10/gosseract"
)

var (
	parseCaptcha = regexp.MustCompile(`<img.+src="([^"]+)"\s*alt="Captcha"`)
)

// DOESN'T ALWAYS WORK... THE ISSUE IS THE CAPTCHA IS A HIT OR MISS
// https://hashkiller.co.uk/sha1-decrypter.aspx much better...
func decrypt(hash string) (pass string, err error) {
	// https://hashkiller.co.uk/sha1-decrypter.aspx POST
	// HEADERS
	/*
		:authority:hashkiller.co.uk
		:method:POST
		:path:/sha1-decrypter.aspx
		:scheme:https
		accept:*/ /*
		accept-encoding:gzip, deflate, br
		accept-language:en-US,en;q=0.8,ca;q=0.6
		cache-control:no-cache
		content-length:550
		content-type:application/x-www-form-urlencoded; charset=UTF-8
		cookie:__cfduid=d5f786273e57f2eb07b7410abcbb0147d1497815718; ASP.NET_SessionId=aufstwskz2x0h225cj25ycts; _ga=GA1.3.822838588.1497815720; _gid=GA1.3.1116298762.1497815720
		dnt:1
		origin:https://hashkiller.co.uk
		referer:https://hashkiller.co.uk/sha1-decrypter.aspx
		user-agent:Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36
		x-microsoftajax:Delta=true
		x-requested-with:XMLHttpRequest
	*/
	// FORM DATA
	/*
		ctl00$ScriptMan1:ctl00$content1$updDecrypt|ctl00$content1$btnSubmit
		ctl00$content1$txtInput:e7442fb10f337d9458fa3fa6d2cd15817ee39b01
		ctl00$content1$txtCaptcha:8!9C(J
		__EVENTTARGET:
		__EVENTARGUMENT:
		__VIEWSTATE:/wEPaA8FDzhkNGI2OGVmZmVmYzU3Y2R8sL7jr8KwOTY8tg5JHepbhBIiaLPYKdCtKueysL8pVg==
		__EVENTVALIDATION:/wEdAAUMkfTke17Y2hg0Mg79vYEWqH4D6ZgR89DUFBTMOlnEF/gi3F50GIGWR1Nab02LintqB32jq7fFWKhYKPHg/KhxevIaiZTm5H7q+peJqAD0HguFCfotnInRR/NItL3q8Jc5PV4Cp2rk1Flhq806fCNW
		__ASYNCPOST:true
		ctl00$content1$btnSubmit:Submit
	*/
	// need to parse captcha!

	res, err := grequests.Get("https://hashkiller.co.uk/sha1-decrypter.aspx", nil)
	if err != nil {
		log.Fatalln(err)
	}
	html := res.String()
	m := parseCaptcha.FindStringSubmatch(html)
	captchaURL := m[1]
	res, err = grequests.Get(fmt.Sprintf("https://hashkiller.co.uk%s", captchaURL), nil)
	if err != nil {
		log.Fatalln(err)
	}
	res.DownloadToFile("cap.JPEG")

	jpegFile, err := os.Open("cap.JPEG")
	if err != nil {
		log.Fatalln(err)
	}
	defer jpegFile.Close()

	imgJ, err := jpeg.Decode(jpegFile)
	if err != nil {
		log.Fatalln(err)
	}
	// grey scale the image
	bounds := imgJ.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(image.Rectangle{Max: image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			var rgba = imgJ.At(x, y)
			gray.Set(x, y, rgba)

		}
	}
	// clean it up some more
	clean := image.NewRGBA64(image.Rectangle{Max: image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			var rgba = gray.At(x, y)
			r, g, b, _ := rgba.RGBA()
			if int(r/257) < 70 && int(g/257) < 70 && int(b/257) < 70 {
				clean.Set(x, y, color.Black)
			} else {
				clean.Set(x, y, color.Transparent)
			}
		}
	}
	newImage := resize.Resize(uint(float64(w)*1.5), uint(float64(h)*1.5), clean, resize.NearestNeighbor)
	var buf bytes.Buffer
	err = png.Encode(&buf, newImage)
	if err != nil {
		log.Fatalln(err)
	}
	os.Remove("cap.JPEG")

	orc, err := gosseract.NewClient()
	if err != nil {
		log.Fatalln(err)
	}
	image, _, err := image.Decode(strings.NewReader(buf.String()))
	if err != nil {
		log.Fatalln(err)
	}
	client := orc.Image(image)
	pass, err = client.Out()
	pass = strings.TrimSpace(pass)
	// ioutil.WriteFile("tst.png", buf.Bytes(), 0666)
	return
}

func main() {
	p, err := decrypt("e7442fb10f337d9458fa3fa6d2cd15817ee39b01") // answer should be oiybtd
	fmt.Println("pas:", p)
	fmt.Println(err)
}
