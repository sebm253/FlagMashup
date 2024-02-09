package utils

import (
	"bytes"
	"flag-mashup/data"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"sort"

	"golang.org/x/exp/maps"
)

func MashupFlags(src, dst string, codeData *data.CodeData) (*bytes.Buffer, error) {
	image.RegisterFormat("jpeg", "jpg", jpeg.Decode, jpeg.DecodeConfig)

	srcBody, err := codeData.FetchFlag(src)
	if err != nil {
		return nil, err
	}
	defer srcBody.Close()
	srcImage, _, err := image.Decode(srcBody)
	if err != nil {
		return nil, err
	}
	_, srcSortedKeys := getProminentImageColors(srcImage)

	dstBody, err := codeData.FetchFlag(dst)
	if err != nil {
		return nil, err
	}
	defer dstBody.Close()
	dstImage, _, err := image.Decode(dstBody)
	if err != nil {
		return nil, err
	}
	dstPixels, dstSortedKeys := getProminentImageColors(dstImage)

	amount := min(3, len(dstSortedKeys), len(srcSortedKeys)) // replace 3 most dominant colors at most
	dstSortedKeys = dstSortedKeys[:amount]

	modifiedImage := image.NewRGBA(dstImage.Bounds())
	draw.Draw(modifiedImage, dstImage.Bounds(), dstImage, image.Point{}, draw.Over)

	for i, key := range dstSortedKeys {
		srcColor := srcSortedKeys[i]
		for _, coords := range dstPixels[key] {
			modifiedImage.Set(coords.X, coords.Y, color.RGBA{R: srcColor.R, G: srcColor.G, B: srcColor.B, A: srcColor.A})
		}
	}

	buf := new(bytes.Buffer)
	return buf, jpeg.Encode(buf, modifiedImage, &jpeg.Options{
		Quality: 100,
	})
}

func getProminentImageColors(image image.Image) (pixelMap map[pixelData][]coordsData, keys []pixelData) { // map of RGBA:slice of all cords
	bounds := image.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixelMap = make(map[pixelData][]coordsData)
	for y := 0; y < height; y++ { // loop through all pixels
		for x := 0; x < width; x++ {
			rgbaPixel := rgbaToPixel(image.At(x, y).RGBA())
			coords := pixelMap[rgbaPixel]
			pixelMap[rgbaPixel] = append(coords, coordsData{X: x, Y: y})
		}
	}
	keys = maps.Keys(pixelMap)
	sort.Slice(keys, func(i, j int) bool {
		return len(pixelMap[keys[i]]) > len(pixelMap[keys[j]])
	})
	return
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) pixelData {
	return pixelData{uint8(r / 257), uint8(g / 257), uint8(b / 257), uint8(a / 257)}
}

type pixelData struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type coordsData struct {
	X int
	Y int
}
