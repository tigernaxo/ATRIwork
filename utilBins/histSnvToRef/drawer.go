package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
	"github.com/tigernaxo/ATRIwork/fileformat"
	"github.com/tigernaxo/ATRIwork/snv"
)

// Draw histogram
// 重要：應該改成每次都resize塞到final Image並釋放記憶體，全部展開會造成out of Memory
func drawHist(gapRatio float64, convasWidth int, ref []byte, fileList []string) *image.RGBA {
	unitConvas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{convasWidth, 1},
	})
	shinkedConvas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{1920, 10},
	})
	gapHeight := int(float64(10) * gapRatio)
	finalConvas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{1920, len(fileList)*10 + gapHeight*(len(fileList)-1)},
	})
	anchorY := 0
	siteMap := mkSiteMap(len(ref), false)
	for i, fa := range fileList {
		_, seq := fileformat.ReadSingleFasta(fa)
		for i := range siteMap {
			siteMap[i] = false
		}
		siteMap = snv.SiteMapUpdate(siteMap, ref, seq)
		if len(siteMap) > convasWidth {
			log.Panicf("[Error] siteMap is larger then convasWidth \n")
		}

		// fill histogram
		for x := 0; x < len(siteMap); x++ {
			if siteMap[x] {
				unitConvas.SetRGBA(x, 0, *snvColor)
			}
		}
		// resize to 1920 width
		shinkedConvas = resize.Resize(1920, 10, unitConvas, resize.NearestNeighbor).(*image.RGBA)
		for x := 0; x < 1920; x++ {
			for y := anchorY; y < anchorY+10; y++ {
				finalConvas.SetRGBA(x, y, shinkedConvas.At(x, 0).(color.RGBA))
			}
		}
		anchorY += 10
		// Draw shinkedConvas on filnalConvas
		if i+1 != len(fileList) {
			// Draw(Skip) gap
			anchorY += gapHeight
		}
	}
	return finalConvas
}
func drawGenome(genomeLength, convasWidth, outHeight, outWidth int) *image.RGBA {
	convas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{convasWidth, 1},
	})
	for x := 0; x < genomeLength; x++ {
		convas.SetRGBA(x, 0, *genomeColor)
	}
	return convas
}
func drawScale(unit, convasWidth int) *image.RGBA {
	// 目標是在1920寬的圖上mainUnit有3px subUnit有1px
	scaleAxisPx := 3
	scaleMainPx := 3
	scaleMainHeight := 10
	// scaleSubPx := 1
	scaleSubHeight := 5
	finalWidth := 1920
	finalHeight := scaleMainHeight + scaleAxisPx + scaleSubHeight
	scaleColor := &color.RGBA{0, 0, 0, 255}

	// subUnit := int(float64(unit) / 10.)
	convas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{finalWidth, finalHeight},
	})
	// Draw Axis
	for x := 0; x < finalWidth; x++ {
		for y := scaleMainHeight; y < scaleMainHeight+scaleAxisPx; y++ {
			convas.SetRGBA(x, y, *scaleColor)
		}
	}
	// Draw main bar
	for x := 0; x < scaleMainPx; x++ {
		for y := 0; y < scaleMainHeight; y++ {
			convas.SetRGBA(x, y, *scaleColor)
		}
	}
	for x := finalWidth - 3; x < finalWidth; x++ {
		for y := 0; y < scaleMainHeight; y++ {
			convas.SetRGBA(x, y, *scaleColor)
		}
	}
	mainStep := int(float64(finalWidth) / (float64(convasWidth) / float64(unit)))
	fmt.Println(mainStep)
	// subStep := int(float64(finalWidth) / (float64(convasWidth) / (float64(unit) / float64(10))))
	for x := mainStep; x < mainStep*(convasWidth/unit); x = x + mainStep {
		for y := 0; y < scaleMainHeight; y++ {
			convas.SetRGBA(x-1, y, *scaleColor)
			convas.SetRGBA(x, y, *scaleColor)
			convas.SetRGBA(x+1, y, *scaleColor)
		}
	}
	// fmt.Println(subStep)
	// for x := subStep; x < subStep*(convasWidth/(unit/10)); x = x + subStep {
	// 	for y := scaleMainHeight + scaleAxisPx; y < finalHeight; y++ {
	// 		convas.SetRGBA(x, y, *scaleColor)
	// 	}
	// }
	return convas
}
func savePng(img image.Image, outPath string) {
	outImage, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outImage.Close()
	png.Encode(outImage, img)
}
func resizeThenSavePng(img image.Image, width, height uint, outPath string) {
	finalImage := resize.Resize(width, height, img, resize.NearestNeighbor)
	savePng(finalImage, outPath)
}
func mkSiteMap(len int, fill bool) []bool {
	siteMap := make([]bool, len)
	for i := range siteMap {
		siteMap[i] = fill
	}
	return siteMap
}
