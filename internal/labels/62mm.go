package labels

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/math/fixed"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func Create62mmItemLabel(text string, qrValue string) ([]byte, error) {

	textSlice := []string{text}
	dpi := float64(300)
	hinting := "none"
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
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	// Draw the text.
	h := font.HintingNone
	switch hinting {
	case "full":
		h = font.HintingFull
	}
	d := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: h,
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

	return writePNG(rotatedImg)
}

func rot90(inputImage *image.RGBA) *image.RGBA {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	bounds.Max.X, bounds.Max.Y = bounds.Max.Y, bounds.Max.X

	outputImage := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			org := inputImage.At(x, y)
			outputImage.Set(height-y, x, org)
		}
	}

	return outputImage
}

type QrWriter struct {
}

func newQRCode(destination string) (image.Image, error) {
	qrc, err := qrcode.New(destination)
	qrc.Dimension()
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)

	bw := bufio.NewWriter(buffer)

	mwc := &MyWriteCloser{bw}

	w := standard.NewWithWriter(mwc,
		standard.WithQRWidth(20),
		// standard.WithCircleShape(),
		// standard.WithBorderWidth(3),
		// standard.WithHalftone("logo.png"),
	)

	if err = qrc.Save(w); err != nil {
		fmt.Printf("could not save image: %v\n", err)
	}

	image, _, err := image.Decode(buffer)

	return image, err
}

func writePNG(img *image.RGBA) ([]byte, error) {
	var buf bytes.Buffer
	// Create a buffered writer that writes to the bytes.Buffer.
	b := bufio.NewWriter(&buf)

	// Encode the image to the buffered writer.
	err := png.Encode(b, img)
	if err != nil {
		return nil, err
	}

	// Flush the buffered writer to ensure all data is written.
	err = b.Flush()
	if err != nil {
		return nil, err
	}

	// Return the bytes.Buffer contents as a byte slice.
	return buf.Bytes(), nil
}
