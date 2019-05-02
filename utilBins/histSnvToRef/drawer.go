package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/nfnt/resize"
	"github.com/tigernaxo/ATRIwork/fileformat"
	"github.com/tigernaxo/ATRIwork/snv"
)

func drawHist(conf *config) *image.RGBA {
	_, ref := fileformat.ReadSingleFasta(conf.refSeq)
	mapLen := len(ref)

	// 只要剩餘的convas佔超過全部的1/5(genome length 1/4)，minUnit就往下調整一個log級數
	minUnit := math.Pow10(int(math.Floor(math.Log10(float64(len(ref) / 4)))))
	gapRatio := conf.gapRatio
	fileList := conf.fileList
	convasWidth := int(math.Ceil(float64(len(ref))/minUnit)) * int(minUnit)

	// pre setting convas
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

	ntMask := []byte{'N', '-'}
	showMap := mkSiteMap(mapLen, true)
	// get showMap
	for _, fa := range fileList {
		_, seq := fileformat.ReadSingleFasta(fa)
		snv.UpdateShowMapByMaskChar(showMap, seq, ntMask)
	}
	// drawing
	anchorY := 0
	snvMap := mkSiteMap(mapLen, false)
	for i, fa := range fileList {
		_, seq := fileformat.ReadSingleFasta(fa)
		setBool(snvMap, false)
		snv.UpdateSNVMap(snvMap, ref, seq)
		if len(seq) > mapLen {
			log.Panicf("[Error] snvMap is smaller then sequence\n")
		}

		// fill histogram
		for x := 0; x < mapLen; x++ {
			if snvMap[x] && showMap[x] {
				for i := maxInt(0, x-int((conf.intensity-1.0)/2.0)); i < minInt(len(snvMap), x+int((conf.intensity-1.0)/2.0)); i++ {
					unitConvas.SetRGBA(i, 0, *snvColor)
				}
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

func minInt(i, j int) int {
	if i < j {
		return i
	}
	return j
}
func maxInt(i, j int) int {
	if i > j {
		return i
	}
	return j
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

// Scale ss
type Scale struct {
	Width  int
	Height int
}

func drawScale(unit, totalLen, convasLen int) *image.RGBA {
	convasWidth := 1920
	scaleColor := &color.RGBA{0, 0, 0, 255}
	mainScale := Scale{3, 10}
	subScale := Scale{1, 5}
	middleLinePix := 3
	convas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{convasWidth + mainScale.Width, middleLinePix + mainScale.Height + subScale.Height},
	})

	recTofill := image.Rectangle{}
	// 繪製mainScale
	recTofill.Min.Y = 0
	recTofill.Max.Y = recTofill.Min.Y + mainScale.Height
	for i := 0; i <= totalLen/unit; i++ {
		recTofill.Min.X = int(float64(i) * float64(unit) * (float64(convasLen) / float64(totalLen)))
		recTofill.Max.X = recTofill.Min.X + mainScale.Width
		fillRect(recTofill, scaleColor, *convas)
	}

	// 繪製subScale
	subUnit := float64(unit) / float64(10)
	recTofill.Min.Y = mainScale.Height + middleLinePix
	recTofill.Max.Y = recTofill.Min.Y + subScale.Height
	for i := 0; i <= int(totalLen/unit)*10; i++ {
		recTofill.Min.X = int(float64(i) * subUnit * (float64(convasLen) / float64(totalLen)))
		recTofill.Max.X = recTofill.Min.X + subScale.Width
		fillRect(recTofill, scaleColor, *convas)
	}
	// 繪製middleLine
	recTofill.Min.X = 0
	recTofill.Max.X = convas.Bounds().Max.X
	recTofill.Min.Y = mainScale.Height
	recTofill.Max.Y = recTofill.Min.Y + middleLinePix
	fillRect(recTofill, scaleColor, *convas)

	return convas
}
func fillRect(rec image.Rectangle, c *color.RGBA, img image.RGBA) {
	for x := rec.Min.X; x < rec.Max.X; x++ {
		for y := rec.Min.Y; y < rec.Max.Y; y++ {
			img.SetRGBA(x, y, *c)
		}
	}
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
func setBool(a []bool, b bool) {
	for i := range a {
		a[i] = b
	}
}
