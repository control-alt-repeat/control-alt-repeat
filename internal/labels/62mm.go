package labels

import (
	"image"
	"image/color"
	"math"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/math/fixed"
)

func Create62mmItemLabel(text string, qrValue string) ([]byte, error) {
	textSlice := []string{text}
	dpi := float64(300)
	size := float64(55)
	spacing := float64(0.5)
	whiteOnBlack := false

	f, err := truetype.Parse(gomonobold.TTF)
	if err != nil {
		return nil, err
	}

	// Draw the background and the guidelines.
	fg, bg := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	if whiteOnBlack {
		fg, bg = image.White, image.Black
		ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	}
	const imgW, imgH = 1109, 696
	rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	d := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: font.HintingNone,
		}),
	}
	y := 10 + int(math.Ceil(size*dpi/72))
	dy := int(math.Ceil(size * spacing * dpi / 72))

	y += dy
	for _, s := range textSlice {
		d.Dot = fixed.P(10, y)
		d.DrawString(s)
		y += dy
	}

	img2, err := newQRCode(qrValue)
	if err != nil {
		return nil, err
	}

	var img1 image.Image = rgba

	sp2 := image.Point{img1.Bounds().Dx(), 0}

	r2 := image.Rectangle{sp2, sp2.Add(img2.Bounds().Size())}

	r := image.Rectangle{image.Point{0, 0}, r2.Max}

	rgba = image.NewRGBA(r)

	draw.Draw(rgba, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, r2, img2, image.Point{0, 0}, draw.Src)

	newHeight := 696

	aspectRatio := float64(rgba.Bounds().Dx()) / float64(rgba.Bounds().Dy())
	newWidth := int(float64(newHeight) * aspectRatio)

	resizedImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), rgba, rgba.Bounds(), draw.Over, nil)

	rotatedImg := rot90(resizedImg)

	return writeRGBAtoPNG(rotatedImg)
}
