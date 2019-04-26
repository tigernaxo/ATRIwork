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
// 重要：應該改成每次都resize塞到final Image，全部展開會out of Memory
func drawHist(histHeight, gapHeight, convasWidth, convasHeight int, ref []byte, fileList []string) *image.RGBA {
	convas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{convasWidth, convasHeight},
	})
	startY := 0
	for _, fa := range fileList {
		fmt.Printf("%s Creating SNV Histogram: %s\n", timeStamp(), fa)

		_, seq := fileformat.ReadSingleFasta(fa)
		siteMap := mkSiteMap(len(ref), false)
		siteMap = snv.SiteMapUpdate(siteMap, ref, seq)

		for x := 0; x < len(siteMap); x++ {
			if siteMap[x] {
				for y := startY; y < startY+histHeight; y++ {
					convas.SetRGBA(x, y, *snvColor)
				}
			}
		}
		startY = startY + histHeight + gapHeight
	}
	return convas
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
