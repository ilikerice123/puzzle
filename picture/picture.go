package picture

import (
	"fmt"
	"image"

	// for allowing image.Decode to understand gif
	_ "image/gif"
	"image/jpeg"

	// for allowing image.Decode to understand png
	_ "image/png"
	"os"
	"path"

	"github.com/disintegration/gift"
)

// TODO: the API for this package is a bit inconsistent:
// - some methods are file-based
// - some return an image
// - some modify the iamge in place
// should try and unify the API

// LoadImage returns an image from a file system
func LoadImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// SaveImage saves an image to the file system with the filename
func SaveImage(filename string, img image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return jpeg.Encode(f, img, nil)
}

// SliceImage slices up an image into a bunch of
func SliceImage(filename string, ySize int, xSize int) ([][]string, error) {
	img, err := LoadImage(filename)
	if err != nil {
		return nil, err
	}
	height, width := normalizeImage(img, ySize, xSize)
	pieceHeight := height / ySize
	pieceLength := width / xSize
	cropRect := image.Rectangle{}

	imageNames := make([][]string, ySize)
	for i := range imageNames {
		imageNames[i] = make([]string, xSize)
	}
	for i := 0; i < ySize; i++ {
		for j := 0; j < xSize; j++ {
			cropRect.Min.Y = i * pieceHeight
			cropRect.Min.X = j * pieceLength
			cropRect.Max.Y = cropRect.Min.Y + pieceHeight
			cropRect.Max.X = cropRect.Min.X + pieceLength

			filter := gift.New(gift.Crop(cropRect))
			dst := image.NewNRGBA(filter.Bounds(img.Bounds()))
			filter.Draw(dst, img)

			ext := path.Ext(filename)
			name := filename[0 : len(filename)-len(ext)]
			fileName := fmt.Sprintf("%s_%d_%d%s", name, i, j, ext)

			SaveImage(fileName, dst)
			imageNames[i][j] = fileName
		}
	}
	return imageNames, nil
}

// DownsizeImage resizes the image for a preview
func DownsizeImage(img image.Image) image.Image {
	filter := gift.New(gift.Resize(200, 0, gift.LanczosResampling))
	dst := image.NewNRGBA(filter.Bounds(img.Bounds()))
	filter.Draw(dst, img)
	return dst
}

// normalizeImage resizes the image so the bounds are a multiple of ySize and xSize
func normalizeImage(img image.Image, ySize int, xSize int) (height int, width int) {
	bounds := img.Bounds()
	height = bounds.Max.Y - bounds.Min.Y
	width = bounds.Max.X - bounds.Min.X
	height = (height / ySize) * ySize
	width = (width / xSize) * xSize
	filter := gift.New(gift.CropToSize(width, height, gift.LeftAnchor))
	dst := image.NewNRGBA(filter.Bounds(img.Bounds()))
	filter.Draw(dst, img)
	img = dst
	return
}
