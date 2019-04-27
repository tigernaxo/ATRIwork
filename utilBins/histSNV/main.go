package main

import (
	"math"

	"github.com/tigernaxo/ATRIwork/fileformat"
)

func main() {
	// 計算各1920位置的強度
	convasWidth := 1920
	refFile := "./holdref.fa"
	currentHist := make([]uint8, 1920)
	_, refSeq := fileformat.ReadSingleFasta(refFile)
	genomeLen := len(refSeq)
	siteMap := make([]bool, genomeLen)
	for i := range siteMap {
		siteMap[i] = false
	}
	unit := float64(genomeLen) / float64(convasWidth)
	span := math.Ceil(unit / 2)
	var sum float64
	for i := 0; i < convasWidth; i++ {
		center := math.Round(float64(genomeLen)*float64(i)/float64(convasWidth) + float64(unit)/2)
		sum = 0
		for j := int(center - span + 1); j < int(center+span-1); j++ {
			if siteMap[j] {
				sum++
			}
		}
		currentHist[i] = uint8(math.Round(255 * (sum / (2*span - 2))))
	}
	// 按照獲得的顏色強弱繪置單元histogram
	// 跳gap
}
