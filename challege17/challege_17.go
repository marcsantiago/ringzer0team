package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"log"
	"regexp"
	"strings"

	"github.com/nfnt/resize"
	"github.com/otiai10/gosseract"

	"../auth"
)

var parseBase = regexp.MustCompile(`<img\s*src="(data[^"]+)`)

func main() {
	c, err := auth.NewSession()
	if err != nil {
		log.Fatalln(err)
	}
	foundFlag := false
	var counter int
	// THE ORC ISN'T PERFECT
	// BUT EVENTUALLY IT WILL MATCH :-)
	for !foundFlag {
		counter++
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
				// see if the pixel is white and get rid of junk
				if r/257 == 255 && g/257 == 255 && b/257 == 255 {
					clean.Set(x, y, color.Black)
				} else {
					clean.Set(x, y, color.Transparent)
				}

			}
		}
		// make the image large so that the OCR works better
		newImage := resize.Resize(uint(w*2), uint(h*2), clean, resize.NearestNeighbor)

		var buf bytes.Buffer
		err = png.Encode(&buf, newImage)
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
		if err != nil {
			log.Fatalln(err)
		}
		if pass == "" {
			continue
		}
		pass = strings.TrimSpace(pass)
		err = c.SubmitAnswer("17", pass)
		if err == nil {
			foundFlag = true
		} else {
			log.Println("interation", counter, err.Error(), "parsed", pass)
		}
	}
}
