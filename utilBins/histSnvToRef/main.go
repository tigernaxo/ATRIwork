package main

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"time"

	"github.com/tigernaxo/ATRIwork/fileformat"
)

var snvColor = &color.RGBA{255, 0, 0, 255}
var genomeColor = &color.RGBA{0, 0, 0, 255}

func main() {
	outHist := "histogram.png"
	outScale := "scale.png"
	outGenome := "genome.png"
	// histHeight := 20

	if len(os.Args) == 1 {
		fmt.Println("Usage example:")
		fmt.Printf("\t%s ./config_file\n", os.Args[0])
	}
	conf := newConfig(os.Args[1])

	// gapHeight := int(float64(histHeight) * conf.gapRatio)
	// convasHeight := histHeight*len(conf.fileList) + (len(conf.fileList)-1)*gapHeight

	_, refSeq := fileformat.ReadSingleFasta(conf.refSeq)
	// 只要剩餘的convas佔超過全部的1/5(genome length 1/4)，minUnit就往下調整一個log級數
	minUnit := math.Pow10(int(math.Floor(math.Log10(float64(len(refSeq) / 4)))))
	convasWidth := int(math.Ceil(float64(len(refSeq))/minUnit)) * int(minUnit)

	// Draw histogram
	convas := drawHist(conf.gapRatio, convasWidth, refSeq, conf.fileList)
	resizeThenSavePng(convas, uint(1920), uint(1080), conf.outDir+"/"+outHist)

	// Draw genome
	convas = drawGenome(len(refSeq), convasWidth, 20, 1920)
	resizeThenSavePng(convas, uint(1920), uint(20), conf.outDir+"/"+outGenome)

	// Draw Scale bar
	// convas = drawScale(int(minUnit), convasWidth)
	convas = drawScale(int(minUnit), convasWidth, 1920)
	savePng(convas, conf.outDir+"/"+outScale)

	fmt.Printf("%s All Done\n", timeStamp())
}
func timeStamp() string {
	t := time.Now()
	return fmt.Sprintf("[%02v:%02v:%02v]", t.Hour(), t.Minute(), t.Second())
}
