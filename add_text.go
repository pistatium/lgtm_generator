package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
)

const (
	assetsFolder  = "assets"
	outputsFolder = "outputs"
)

func processImages() error {
	files, err := os.ReadDir(assetsFolder)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && isImage(file.Name()) {
			if err := processImage(file.Name()); err != nil {
				return err
			}
		}
	}

	return nil
}

func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

func processImage(filename string) error {
	inputPath := filepath.Join(assetsFolder, filename)
	imgFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}

	out := image.NewRGBA(img.Bounds())
	draw.Draw(out, out.Bounds(), img, image.Point{0, 0}, draw.Src)

	// Initialize the graphics context
	dc := gg.NewContextForRGBA(out)
	width := float64(img.Bounds().Dx())
	height := float64(img.Bounds().Dy())

	// Set the font, size, and text color
	if err := dc.LoadFontFace("./assets/font/Sora-ExtraBold.ttf", height/3); err != nil {
		return err
	}

	// Draw the text shadow
	sw, sh := dc.MeasureString("LGTM")
	x := (width - sw) / 2
	y := (height + sh*0) / 2

	// Black for shadow
	dc.SetRGBA(0.9, 0.9, 0.9, 0.5)
	for dx := 0.0; dx <= 1.0; dx += 1 {
		for dy := 0.0; dy <= 1.0; dy += 1 {
			dc.DrawString("LGTM", x+dx, y+dy)
		}
	}
	// white
	dc.SetRGBA(1, 1, 1, 0.9)
	dc.DrawString("LGTM", x, y)

	outputPath := filepath.Join(outputsFolder, strings.TrimSuffix(filename, filepath.Ext(filename))+".png")
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err := png.Encode(outFile, out); err != nil {
		return err
	}

	fmt.Printf("Processed %s and saved to %s\n", inputPath, outputPath)
	return nil
}

func main() {
	if err := processImages(); err != nil {
		fmt.Println("Error:", err)
	}
}
