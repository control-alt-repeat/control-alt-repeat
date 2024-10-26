package labels

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/opentype"
)

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

func writeRGBAtoPNG(img *image.RGBA) ([]byte, error) {
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

func writeImageToPNG(img image.Image) ([]byte, error) {
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
		standard.WithQRWidth(10),
		// standard.WithQRWidth(15),
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

func loadFontFace(size float64) font.Face {

	f, err := opentype.Parse(gomonobold.TTF)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size, // Set the desired font size here
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("failed to create font face: %v", err)
	}

	return face
}
