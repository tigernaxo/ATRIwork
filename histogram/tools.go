package histogram

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/nfnt/resize"
)

// ImageCombiner 批次結合影像，方便附在phylogenetic tree後面對照
// EdgeNoZero的話最小會有1pixel
type ImageCombiner struct {
	Images          []image.Image
	FinalSize       image.Rectangle
	GapToHistoRatio float64
	GapColor        color.Color
	GapNoZero       bool
}

// CombineImage take images and draw them in one image, then resize
// 要如何在resize之後控制EdgeLine = 1??
// 計算出edgeLine最後是1px的畫要佔去多少px然後四捨五入
// 建立convas: w, h := image.w, sum.image.h + edgeLine.h (edgeLine的數量是image數量-1)
// 逐一放入 image/edgeLine/image/edgeLIne...
// 最後resize convas
func (ic *ImageCombiner) CombineImage() (finalImage image.Image) {
	// 計算全部image的總長，最寬
	var imgMaxX, imgTotalY int
	for _, img := range ic.Images {
		if imgMaxX < img.Bounds().Max.X {
			imgMaxX = img.Bounds().Max.X
		}
		imgTotalY += img.Bounds().Max.Y
	}
	// 計算Gap，並製造Gap Image
	// 假設Gap高度x，縮放後最小高度為1
	// 算式x*ic.FinalSize.Max.Y/(x*len(ic.Images)+imgTotalY)=1
	// x*ic.FinalSize.Max.Y=x*len(ic.Images)+imgTotalY
	// x=imgTotalY/(ic.FinalSize.Max.Y-len(ic.Images)) ...無條件進位
	var gapHeight int
	if ic.GapNoZero {
		gapHeight = int(math.Ceil(float64(imgTotalY) / float64((ic.FinalSize.Max.Y - len(ic.Images)))))
	} else {
		gapHeight = int(math.Ceil(float64(ic.Images[0].Bounds().Max.Y) * ic.GapToHistoRatio))
	}

	gapImage := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{imgMaxX, gapHeight},
		})

	// 建立convas
	convasY := gapHeight*len(ic.Images) + imgTotalY
	convas := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{imgMaxX, convasY},
		})

	// 從0點開始逐一貼上images和Gap
	nowY := 0
	nowPoint := image.Point{0, nowY}
	var convasPaintArea image.Rectangle
	for i, img := range ic.Images {

		convasPaintArea = image.Rectangle{
			nowPoint,
			nowPoint.Add(img.Bounds().Max),
		}
		draw.Draw(convas, convasPaintArea, img, img.Bounds().Min, draw.Src)
		nowY += img.Bounds().Max.Y
		nowPoint = image.Point{0, nowY}

		// drawGap ,return if last image combined
		if i >= len(ic.Images)-1 {
			return
		}
		convasPaintArea = image.Rectangle{
			nowPoint,
			nowPoint.Add(gapImage.Bounds().Max),
		}
		draw.Draw(convas, convasPaintArea, gapImage, gapImage.Bounds().Min, draw.Src)
		nowY += gapHeight
		nowPoint = image.Point{0, nowY}
	}
	return resize.Resize(
		uint(ic.FinalSize.Max.Y),
		uint(ic.FinalSize.Max.X),
		convas,
		resize.NearestNeighbor,
	)
}
