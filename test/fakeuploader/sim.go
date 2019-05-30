package fakeuploader

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"net/http"
	"time"

	"github.com/graeme-hill/gnet/sys/gnet"
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
		X: topLeft.X + rand.Intn(w-topLeft.X) + 1,
		Y: topLeft.Y + rand.Intn(h-topLeft.Y) + 1,
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

func randomImageSize() (int, int) {
	// only use portrait 20% of the time
	portrait := rand.Intn(5) == 0
	if portrait {
		return 768, 1024
	}
	return 1024, 768
}

func randomImage() (*bytes.Buffer, error) {
	// make image
	width, height := randomImageSize()
	topLeft := image.Point{X: 0, Y: 0}
	bottomRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{topLeft, bottomRight})

	drawBackground(img)
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

	return &pngBuffer, nil
}

func continuouslyUpload(ctx context.Context, over chan<- error, apiAddress string) {
	for {
		select {
		case <-ctx.Done():
			over <- nil
			return
		case <-time.After(10 * time.Second):
			fmt.Println("fake client uploading image")
			png, err := randomImage()
			if err != nil {
				over <- err
				return
			}
			_, err = http.Post(apiAddress, "image/png", png)
			if err != nil {
				over <- err
				return
			}
		}
	}
}

func Run(ctx context.Context, apiAddress string) gnet.Service {
	rand.Seed(1)
	over := make(chan error)
	go continuouslyUpload(ctx, over, apiAddress)
	return over
}
