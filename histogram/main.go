package histogram

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/nfnt/resize"
	"github.com/tigernaxo/ATRIwork/fileformat"
)

// accept: width vcf-file maxLength
func main() {

	vcfFile := os.Args[1]
	width, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Panic(err)
	}
	minUnit := 1000000
	genomeLength := 4827641
	maxLength := int(math.Ceil(float64(genomeLength) / float64(minUnit)))
	drawColor := color.NRGBA{uint8(255), 0, 0, 255}
	bgColor := color.NRGBA{218, 236, 198, 255}
	redundantColor := color.NRGBA{255, 255, 255, 255}

	snps := fileformat.SiteFromVCF(vcfFile)
	file, err := os.Create(vcfFile + ".jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	rgba := image.NewRGBA(image.Rect(0, 0, maxLength, 1))
	for x := 0; x < maxLength; x++ {
		rgba.Set(x, 0, bgColor)
	}

	for x := 0; x < len(snps); x++ {
		rgba.Set(snps[x], 0, drawColor)
		if snps[x]+width < maxLength && snps[x]-width > 0 {
			for i := snps[x] - width; i <= snps[x]+width; i++ {
				rgba.Set(i, 0, drawColor)
			}
		} else if snps[x]+width > maxLength && snps[x]-width > 0 {
			for i := snps[x] - width; i <= maxLength; i++ {
				rgba.Set(i, 0, drawColor)
			}
		} else if snps[x]-width < 0 && snps[x]+width < maxLength {
			for i := 0; i < snps[x]+width; i++ {
				rgba.Set(i, 0, drawColor)
			}
		}
	}
	if maxLength > genomeLength {
		for i := genomeLength + 1; i <= maxLength; i++ {
			rgba.Set(i, 0, redundantColor)
		}

	}
	resizedImg := resize.Resize(1920, 20, rgba, resize.NearestNeighbor)

	// fmt.Println(rgba.At(400, 100))     //{144 100 0 255}
	// fmt.Println(rgba.Bounds())         //(0,0)-(500,200)
	// fmt.Println(rgba.Opaque())         //true，其完全透明
	// fmt.Println(rgba.PixOffset(1, 1))  //2004
	// fmt.Println(rgba.Stride)           //2000
	jpeg.Encode(file, resizedImg, nil) //将image信息存入文件中
}
