package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/tigernaxo/ATRIwork/fileformat"
	"github.com/tigernaxo/ATRIwork/histogram"
	"github.com/tigernaxo/ATRIwork/snv"
)

// Input:
// 1. Intensity
// 2. ref fasta
// 3. ...seqs

// Fixed:
// 4.color
// 5.bgcolor
// 6.min Unit ( ex:1000000, to decide convasLength )
// 7.out dimension
// 8.out name

func main() {
	intensity, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Panic(err)
	}
	ref := os.Args[2]
	_, refSeq := fileformat.ReadSingleFasta(ref)

	seqs := os.Args[3:]
	// 只要剩餘的convas佔超過全部的1/5(genome length 1/4)，minUnit就往下調整一個log級數
	minUnit := math.Pow10(int(math.Floor(math.Log(float64(len(refSeq) / 4)))))
	convasWidth := int(math.Ceil(float64(len(refSeq)) / minUnit))

	for _, fa := range seqs {
		// 計算site map
		_, seq := fileformat.ReadSingleFasta(fa)
		siteMap := make([]bool, len(refSeq))
		for i := range siteMap {
			siteMap[i] = false
		}
		siteMap = snv.SiteMapUpdate(siteMap, refSeq, seq)

		p := &histogram.PlotSites{
			SitesMap:     siteMap,
			RefLen:       len(refSeq),
			ConvasLen:    convasWidth,
			Color:        &color.RGBA{255, 0, 0, 255},
			Bgcolor:      &color.RGBA{0, 255, 0, 0},
			Intensity:    float64(intensity),
			OutName:      fa + ".jpg",
			OutDimension: image.Rectangle{image.Point{0, 0}, image.Point{1920, 20}},
		}
		p.PlotSites()
	}
}
