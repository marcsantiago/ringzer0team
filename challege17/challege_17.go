package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/otiai10/gosseract"

	"../auth"
)

var parseBase = regexp.MustCompile(`<img\s*src="(data[^"]+)`)

func main() {
	c, err := auth.NewSession()
	if err != nil {
		log.Fatalln(err)
	}

	res, err := c.Session.Get("https://ringzer0team.com/challenges/17", nil)
	if err != nil {
		log.Fatalln(err)
	}
	html := res.String()
	m := parseBase.FindStringSubmatch(html)
	base := strings.Replace(m[1], "data:image/png;base64,", "", -1)

	imageReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base))
	pngImage, _, err := image.Decode(imageReader)
	if err != nil {
		log.Fatalln(err)
	}

	bounds := pngImage.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	clean := image.NewRGBA64(image.Rectangle{Max: image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			var rgba = pngImage.At(x, y)
			r, g, b, _ := rgba.RGBA()
			// see if the pixel is white
			if r/257 == 255 && g/257 == 255 && b/257 == 255 {
				clean.Set(x, y, color.Black)
			} else {
				clean.Set(x, y, color.Transparent)
			}

		}
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, clean)
	if err != nil {
		log.Fatalln(err)
	}

	orc, err := gosseract.NewClient()
	if err != nil {
		log.Fatalln(err)
	}
	image, _, err := image.Decode(strings.NewReader(buf.String()))
	if err != nil {
		log.Fatalln(err)
	}
	client := orc.Image(image)
	pass, err := client.Out()
	fmt.Println("pass", pass, err)
	ioutil.WriteFile("test.png", buf.Bytes(), 0666)

}

// w 192
// h 55
// res 72
