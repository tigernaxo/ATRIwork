package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
	"github.com/tigernaxo/ATRIwork/snv"

	"github.com/tigernaxo/ATRIwork/fileformat"
)

func main() {
	// Using ref length as map length
	_, r := fileformat.ReadSingleFasta(os.Args[1])
	mapLen := len(r)
	r = nil

	fastas := os.Args[2:]
	maskChar := []byte{'N', '-'}
	showMap := make([]bool, mapLen)
	setBool(showMap, true)
	snvMap := make([]bool, mapLen)
	setBool(snvMap, false)
	maskedShowMap := make([]bool, mapLen)

	var ref, seq []byte
	for i, fa := range fastas {
		_, seq = fileformat.ReadSingleFasta(fa)
		snv.UpdateShowMapByMaskChar(showMap, seq, maskChar)
		switch i {
		case 0:
			ref = seq
		default:
			snv.UpdateSNVMap(snvMap, ref, seq)
		}
	}
	andBoolArray(showMap, snvMap, maskedShowMap)
	//
	strip := &SnvStrip{
		NtColorMap: map[byte]*color.RGBA{
			'A': &color.RGBA{255, 0, 0, 255},
			'T': &color.RGBA{0, 255, 0, 255},
			'C': &color.RGBA{0, 0, 255, 255},
			'G': &color.RGBA{0, 0, 0, 255},
		},
	}
	//
	eachNtNum := countBool(maskedShowMap, true)
	fmt.Printf("true number: snvMap(%d), showMap(%d)\n", countBool(snvMap, true), countBool(showMap, true))
	ntArr := make([]byte, eachNtNum)
	var ntArrIdx int
	img := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{eachNtNum, 1},
	})
	var shrinkedImg *image.RGBA
	histHeight, gapHeight := 10, 2
	finalWidth := 1920
	convasHeight := len(fastas)*histHeight + (len(fastas)-1)*gapHeight
	fmt.Printf("convas height: %d\n", convasHeight)
	fmt.Printf("seq number: %d\n", len(fastas))
	fmt.Printf("snv to shaw number: %d\n", eachNtNum)
	convas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{finalWidth, convasHeight},
	})
	drawStartY := 0
	for _, fa := range fastas {
		ntArrIdx = 0
		_, seq = fileformat.ReadSingleFasta(fa)
		for i, nt := range seq {
			if maskedShowMap[i] {
				ntArr[ntArrIdx] = nt
				ntArrIdx++
			}
		}
		strip.CleanThenDraw(ntArr, img)
		shrinkedImg = resize.Resize(uint(finalWidth), 1, img, resize.NearestNeighbor).(*image.RGBA)
		for x := 0; x < shrinkedImg.Bounds().Max.X; x++ {
			for y := drawStartY; y < drawStartY+histHeight; y++ {
				convas.SetRGBA(x, y, shrinkedImg.At(x, 0).(color.RGBA))
			}
		}
		drawStartY = drawStartY + histHeight + gapHeight
	}
	savePng(convas, "out.png")
}

type NtColor struct {
	Nt      byte
	NtColor *color.RGBA
}
type SnvStrip struct {
	NtColorMap map[byte]*color.RGBA
}

func (s *SnvStrip) CleanThenDraw(nt []byte, convas *image.RGBA) {
	// Check convas widthcolorConfig[nt[x]]
	// if len(nt) != convas.Bounds().Max.X-1 {
	// 	log.Panic("[Error] Len(nt to draw) != convas width")
	// 	log.Panicf("nt number: %d, convas bounds: %d\n", len(nt), convas.Bounds().Max.X)
	// }
	// Clear convas
	bgColor := &color.RGBA{0, 0, 0, 0}
	for x := 0; x < convas.Bounds().Max.X; x++ {
		for y := 0; y < convas.Bounds().Max.Y; y++ {
			convas.SetRGBA(x, y, *bgColor)
		}
	}
	// Draw
	var currentColor *color.RGBA
	for x := 0; x < len(nt); x++ {
		for t, c := range s.NtColorMap {
			if snv.IsEqualAlphabet(t, nt[x]) {
				currentColor = c
			}
		}
		for y := 0; y < convas.Bounds().Max.Y; y++ {
			convas.SetRGBA(x, y, *currentColor)
		}
	}
}
func countBool(a []bool, b bool) int {
	count := 0
	for _, element := range a {
		if b == element {
			count++
		}
	}
	return count
}
func andBoolArray(aArr []bool, bArr []bool, outArr []bool) {
	if len(aArr) != len(bArr) || len(bArr) != len(outArr) {
		log.Panicf("[Error] not equal length of array.\n")
	}
	for i := 0; i < len(aArr); i++ {
		outArr[i] = aArr[i] && bArr[i]
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
func setBool(a []bool, b bool) {
	for i := range a {
		a[i] = b
	}
}
