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

// 重要：應該改成每次都resize塞到final Image，全部展開會out of Memory
func main() {
	outHist := "histogram.png"
	outScale := "scale.png"
	outGenome := "genome.png"
	histHeight := 20

	conf := newConfig(os.Args[1])

	gapHeight := int(float64(histHeight) * conf.gapRatio)
	convasHeight := histHeight*len(conf.fileList) + (len(conf.fileList)-1)*gapHeight

	_, refSeq := fileformat.ReadSingleFasta(conf.refSeq)
	// 只要剩餘的convas佔超過全部的1/5(genome length 1/4)，minUnit就往下調整一個log級數
	minUnit := math.Pow10(int(math.Floor(math.Log10(float64(len(refSeq) / 4)))))
	convasWidth := int(math.Ceil(float64(len(refSeq))/minUnit)) * int(minUnit)

	// Draw histogram
	convas := drawHist(histHeight, gapHeight, convasWidth, convasHeight, refSeq, conf.fileList)
	resizeThenSavePng(convas, uint(1920), uint(1080), conf.outDir+"/"+outHist)

	// Draw genome
	convas = drawGenome(len(refSeq), convasWidth, 20, 1920)
	resizeThenSavePng(convas, uint(1920), uint(20), conf.outDir+"/"+outGenome)

	// Draw Scale bar
	convas = drawScale(int(minUnit), convasWidth)
	savePng(convas, conf.outDir+"/"+outScale)

	fmt.Printf("%s All Done\n", timeStamp())
}
func timeStamp() string {
	t := time.Now()
	return fmt.Sprintf("[%02v:%02v:%02v]", t.Hour(), t.Minute(), t.Second())
}
