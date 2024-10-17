package labels

import (
	"image"
	"image/color"
	"image/draw"
	"strings"

	"golang.org/x/image/font"

	"golang.org/x/image/math/fixed"
)

func Create102x152mmItemLabel(warehouseID string, itemDescription string, qrValue string) ([]byte, error) {
	width := 1660
	height := 1164
	// marginTop := 300
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	bgColor := color.RGBA{255, 255, 255, 255} // White background
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Draw the QR code on the right-hand side
	qrCode, err := newQRCode(qrValue) // QR Code content
	if err != nil {
		return nil, err
	}
	qrCodeSize := 600                         // Size of QR code
	qrPos := image.Pt(width-qrCodeSize-20, 0) // Position with margin
	draw.Draw(img, image.Rect(qrPos.X, qrPos.Y, qrPos.X+qrCodeSize, qrPos.Y+qrCodeSize), qrCode, image.Point{}, draw.Over)

	drawTextBox(img, 100, 500, 800, itemDescription, loadFontFace(50))

	// Draw text in the bottom left corner
	addLabel(img, 20, 200, warehouseID, loadFontFace(250))

	// return writeImageToPNG(rot90(img))
	return writeImageToPNG(img)
}

// addText adds text to the image at specified coordinates
func addLabel(img *image.RGBA, x, y int, label string, face font.Face) {
	col := color.RGBA{0, 0, 0, 255} // Black text color
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(label)
}

func drawTextBox(img *image.RGBA, x, y, maxWidth int, text string, face font.Face) {
	col := color.RGBA{0, 0, 0, 255} // Black text color
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}

	// Split the text into words
	words := strings.Fields(text)
	spaceWidth := drawer.MeasureString(" ")

	lineHeight := face.Metrics().Height.Ceil()
	startX := point.X

	for _, word := range words {
		wordWidth := drawer.MeasureString(word)

		// Move to next line if word exceeds the max width
		if point.X+wordWidth >= fixed.I(x+maxWidth) {
			point.X = startX
			point.Y += fixed.I(lineHeight)
			drawer.Dot = point
		}

		drawer.DrawString(word)
		point.X += wordWidth + spaceWidth // Move position for next word
		drawer.Dot = point
	}
}
