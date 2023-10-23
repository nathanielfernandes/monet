package main

import (
	"bytes"
	"image"
	"image/png"

	"github.com/nfnt/resize"
)

const WIDTH = 320
const HEIGHT = 320
const PIXEL_SIZE = 16
const W = WIDTH / PIXEL_SIZE
const H = HEIGHT / PIXEL_SIZE

func PayloadToImage(payload string) (*bytes.Buffer, error) {
	im := image.NewRGBA(image.Rect(0, 0, W, H))

	// iterate over the runes
	for i, c := range payload {
		// lookup the color
		color := Lookup(c)

		// calculate the x and y coordinates
		x := i % W
		y := i / W

		im.Set(x, y, color)
	}

	// resize the image to the desired size
	rs := resize.Resize(WIDTH, HEIGHT, im, resize.NearestNeighbor)

	return encodePng(rs)
}

func encodePng(im image.Image) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, im)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
