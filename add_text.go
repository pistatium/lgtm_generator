package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func main() {
	// 画像を読み込む
	file, err := os.Open("./assets/cat_run.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// 画像の中央に LGTM の文字を描画する
	const lgtmText = "LGTM"
	textWidth := 800
	textHeight := 300
	textImg := image.NewRGBA(image.Rect(0, 0, textWidth, textHeight))
	draw.Draw(textImg, textImg.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.Point{}, draw.Src)
	drawString(textImg, lgtmText, 0, 0, color.White, "Sora-ExtraBold")
	x := (img.Bounds().Max.X - textWidth) / 2
	y := (img.Bounds().Max.Y - textHeight) / 2
	draw.Draw(img.(draw.Image), image.Rect(x, y, x+textWidth, y+textHeight), textImg, image.Point{}, draw.Over)

	// 出力する
	out, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		log.Fatal(err)
	}
}
func drawString(dst draw.Image, s string, x, y int, c color.Color, fontName string) {
	f, err := loadFont(fontName)
	if err != nil {
		log.Fatal(err)
	}

	// draw.DrawTextを使用してテキストを描画する
	face := truetype.NewFace(f, &truetype.Options{
		Size:    220,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(c),
		Face: face,
	}
	d.Dot = fixed.P(x, y+face.Metrics().Ascent.Ceil())

	d.DrawString(s)
}

func loadFont(name string) (*truetype.Font, error) {
	f, err := os.Open(filepath.Join("assets", "font", name+".ttf"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return truetype.Parse(b)
}

func getAverageColor(img image.Image) color.RGBA {
	var r, g, b, a uint32
	var count uint32
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			r += uint32(c.R)
			g += uint32(c.G)
			b += uint32(c.B)
			a += uint32(c.A)
			count++
		}
	}
	if count == 0 {
		return color.RGBA{0, 0, 0, 255}
	}
	r /= count
	g /= count
	b /= count
	a /= count
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
