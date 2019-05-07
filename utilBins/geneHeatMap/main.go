package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"sort"

	"github.com/tigernaxo/ATRIwork/fileformat"
)

func main() {
	feature := "CDS"
	gffs := os.Args[1:]
	histHigh, gapHeight := 10, 2
	heatColor := &color.RGBA{255, 0, 0, 255}
	coreColor := &color.RGBA{0, 0, 0, 255}

	// 先union gff gene，一邊累積各gene在sample出現數量作為排序依據
	genes := make(geneCountMap)
	for _, gff := range gffs {
		fileformat.FeatureGeneCountAccumulate(genes, feature, gff)
	}
	// 逐一
	sortedGeneArr := genes.sortToArr()

	fmt.Printf("length of geneCountMap: %d\n", len(genes))
	fmt.Printf("length of sortedGeneArr: %d\n", len(sortedGeneArr))
	startY := 0
	convas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{len(sortedGeneArr), len(gffs)*histHigh + (len(gffs)-1)*gapHeight},
	})
	fmt.Printf("gff number: %d\n", len(gffs))
	fmt.Printf("convas Height: %d\n", convas.Bounds().Max.Y)
	fmt.Printf("convas Width: %d\n", convas.Bounds().Max.X)
	for _, g := range sortedGeneArr {
		fmt.Printf("%s\t%d\n", g.Name, g.Count)
	}
	for _, gff := range gffs {
		heatMap := make([]bool, len(sortedGeneArr))
		setBool(heatMap, false)
		currentGeneMap := make(geneCountMap)
		// 取得heat map
		fileformat.FeatureGeneCountAccumulate(currentGeneMap, feature, gff)
		for i, g := range sortedGeneArr {
			if _, ok := currentGeneMap[g.Name]; ok {
				heatMap[i] = true
			} else {
				heatMap[i] = false
			}
		}
		i := 0
		for _, b := range heatMap {
			if b {
				i++
			}
		}
		fmt.Printf("Union: %d, Find: %d, HeatMap: %d/%d, %s\n", len(sortedGeneArr), len(currentGeneMap), i, len(heatMap), gff)
		// 畫圖
		for x := 0; x < len(sortedGeneArr); x++ {
			if heatMap[x] {
				for y := startY; y < startY+histHigh; y++ {
					switch sortedGeneArr[x].Count {
					case len(gffs):
						convas.SetRGBA(x, y, *coreColor)
					default:
						convas.SetRGBA(x, y, *heatColor)
					}
				}
			}
		}
		startY = startY + histHigh + gapHeight
	}
	savePng(convas, "out.png")
}

type geneCountUnit struct {
	Name  string
	Count int
}
type geneCountArr []geneCountUnit

func (p geneCountArr) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p geneCountArr) Len() int           { return len(p) }
func (p geneCountArr) Less(i, j int) bool { return p[i].Count < p[j].Count }

type geneCountMap map[string]int

func (a geneCountMap) sortToArr() geneCountArr {
	arr := make(geneCountArr, len(a))
	i := 0
	for k, v := range a {
		arr[i] = geneCountUnit{k, v}
		i++
	}
	sort.Sort(arr)
	return arr
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
