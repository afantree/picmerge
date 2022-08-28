package main

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"
	"time"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("please input more pictures")
		os.Exit(-1)
	}
	picPaths := os.Args[1:]
	pics := make([]image.Image, len(picPaths))
	for i, path := range picPaths {
		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		img, formatName, err := image.Decode(f)
		if err != nil {
			panic(err)
		}
		fmt.Printf("read picture %s（%s）\n", path, formatName)
		pics[i] = img
	}
	// new image
	width := pics[0].Bounds().Max.X
	high := pics[0].Bounds().Max.Y
	newImg := image.NewRGBA(image.Rect(0, 0, width*len(picPaths), high))
	for i, pic := range pics {
		rect := image.Rect(width*i, 0, newImg.Bounds().Max.X, newImg.Bounds().Max.Y)
		draw.Draw(newImg, rect, pic, image.Pt(0, 0), draw.Over)
	}

	// check exits config
	var conf Config
	if file, err := os.Open("picmerge.toml"); err == nil {
		if b, errRead := io.ReadAll(file); errRead == nil {
			if _, errDecode := toml.Decode(string(b), &conf); errDecode == nil {
				fmt.Println("decode config")
			}
		}
	}
	var outImage image.Image
	if conf.isEmpty() {
		outImage = newImg
	} else { // subImage
		outImage = newImg.SubImage(image.Rect(
			newImg.Bounds().Min.X+conf.Left,
			newImg.Bounds().Min.Y+conf.Top,
			newImg.Bounds().Max.X-conf.Right,
			newImg.Bounds().Max.Y-conf.Bottom,
		))
	}

	// output image
	outFile, err := os.Create(fmt.Sprintf("%d.jpg", time.Now().Unix()))
	defer outFile.Close()
	if err != nil {
		panic(err)
	}
	b := bufio.NewWriter(outFile)
	err = jpeg.Encode(b, outImage, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		panic(err)
	}
	err = b.Flush()
	if err != nil {
		panic(err)
	}
}

type Config struct {
	Top    int `toml:"top"`
	Bottom int `toml:"bottom"`
	Left   int `toml:"left"`
	Right  int `toml:"right"`
}

func (c Config) isEmpty() bool {
	return c.Top == 0 && c.Bottom == 0 && c.Left == 0 && c.Right == 0
}
