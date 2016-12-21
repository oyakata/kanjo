package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
)

func NewImage() *image.RGBA {
	var w, h int = 280, 240
	m := image.NewRGBA(image.Rect(0, 0, w, h))

	pin := func(x, y int) uint8 {
		var r float64 = 80
		var X, Y float64 = 140, 240

		xx, yy := X - float64(x), Y - float64(y)
		if math.Sqrt(xx*xx + yy*yy) / r > 1 {
			return 200
		} else {
			return 255
		}
	}

	pin2 := func(x, y int) uint8 {
		var r float64 = 50
		var X, Y float64 = 140, 180

		xx, yy := X - float64(x), Y - float64(y)
		if math.Sqrt(xx*xx + yy*yy) / r > 1 {
			return 200
		} else {
			return 255
		}
	}

	mikan := func(x, y int) bool {
		var r float64 = 20
		var X, Y float64 = 140, 120

		xx, yy := X - float64(x), Y - float64(y)
		if math.Sqrt(xx*xx + yy*yy) / r > 1 {
			return true
		} else {
			return false
		}
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			var c color.RGBA
			if mikan(x, y) {
				var value uint8 = uint8(math.Max(float64(pin(x, y)), float64(pin2(x, y))))
				c = color.RGBA{
					value,
					value,
					value,
					255,
				}
			} else {
				c = color.RGBA{
					255,
					102,
					0,
					255,
				}
			}
			m.Set(x, y, c)
		}
	}
	return m
}

func main() {
	img := NewImage()

	tmp, err := ioutil.TempFile("", "SampleImage-")
	if err != nil {
		fmt.Print(err)
		return
	}
	defer os.Remove(tmp.Name())
	png.Encode(tmp, img)
	exec.Command("eog", tmp.Name()).Run()
}
