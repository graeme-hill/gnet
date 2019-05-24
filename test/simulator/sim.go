package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"bufio"
	"bytes"
	"io/ioutil"
	"math/rand"
)

func randomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}

func randomRect(img *image.RGBA) image.Rectangle {
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	topLeft := image.Point{
		X: rand.Intn(w - 1),
		Y: rand.Intn(h - 1),
	}
	bottomRight := image.Point{
		X: topLeft.X + rand.Intn(w - topLeft.X) + 1,
		Y: topLeft.Y + rand.Intn(h - topLeft.Y) + 1,
	}
	
	return image.Rectangle{
		Min: topLeft,
		Max: bottomRight,
	}
}

func drawBackground(img *image.RGBA) {
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	clr := randomColor()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, clr)
		}
	}
}

func drawRect(img *image.RGBA) {
	rect := randomRect(img)
	topLeft := rect.Min
	bottomRight := rect.Max
	clr := randomColor()
	for x := topLeft.X; x < bottomRight.X; x++ {
		for y := topLeft.Y; y < bottomRight.Y; y++ {
			img.Set(x, y, clr)
		}
	}
}

func randomImage() ([]byte, error) {
	rand.Seed(44)

	// make image
	width := 1024
	height := 768
	topLeft := image.Point{X: 0, Y: 0}
	bottomRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{topLeft, bottomRight})

	drawBackground(img)
	drawRect(img)
	drawRect(img)
	drawRect(img)
	drawRect(img)
	drawRect(img)
	drawRect(img)
	drawRect(img)
	drawRect(img)
	drawRect(img)

	// encode as png
	var pngBuffer bytes.Buffer
	writer := bufio.NewWriter(&pngBuffer)
	err := png.Encode(writer, img)
	if err != nil {
		return nil, err
	}

	return pngBuffer.Bytes(), nil
}

func main() {
	png, err := randomImage()
	if err != nil {
		fmt.Printf("oh nooooo, failed to make image b/c '%v'\n", err)
		return
	}

	err = ioutil.WriteFile("C:\\users\\graeme\\test3.png", png, 0644)
	if err != nil {
		fmt.Printf("oh nooooo, failed to write file b/c '%v'\n", err)
		return
	}

	fmt.Println("done")
}
