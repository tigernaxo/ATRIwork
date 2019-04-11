package histogram

import (
	"image"
	"image/color"
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

// CombinePng take images and draw them in one image, then resize
// 要如何在resize之後控制EdgeLine = 1??
// 計算出edgeLine最後是1px的畫要佔去多少px然後四捨五入
// 建立convas: w, h := image.w, sum.image.h + edgeLine.h (edgeLine的數量是image數量-1)
// 逐一放入 image/edgeLine/image/edgeLIne...
// 最後resize convas
func (ic *ImageCombiner) CombinePng() (image image.Image) {
	return ic.Images[0]
}
