package histogram

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/nfnt/resize"
)

// PlotSites contain configs for site histogram
type PlotSites struct {
	SitesMap     []bool
	RefLen       int
	ConvasLen    int
	Color        color.Color
	Bgcolor      color.Color
	Tailcolor    color.Color
	Intensity    float64
	OutName      string
	OutDimension image.Rectangle
}

// PlotSites plot histogram according to PlotSites config information
func (p *PlotSites) PlotSites() {

	if p.RefLen > p.ConvasLen {
		log.Panic("Error: reference length is larger than convas size")
	} else if len(p.SitesMap) != p.RefLen {
		log.Panic("Error: reference length is not equal to length of siteMap")
	}

	// Create png file to output
	file, err := os.Create(p.OutName + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create an image and draw bgcolor in genome range
	rgba := image.NewRGBA(image.Rect(0, 0, p.ConvasLen, 1))

	for i, b := range p.SitesMap {
		if b {
			start := int(math.Round(math.Max(float64(i)-p.Intensity/2, 0)))
			end := int(math.Round(math.Min(float64(i)+p.Intensity/2, float64(len(p.SitesMap)))))
			for j := start; j <= end; j++ {
				rgba.Set(i, 0, p.Color)
			}
		} else {
			rgba.Set(i, 0, p.Bgcolor)
		}
	}
	for i := len(p.SitesMap); i < p.ConvasLen; i++ {
		rgba.Set(i, 0, p.Tailcolor)
	}

	// Resize and save台北內湖行愛理貨中心
	resizedImg := resize.Resize(uint(p.OutDimension.Max.Y), uint(p.OutDimension.Max.X), rgba, resize.NearestNeighbor)
	png.Encode(file, resizedImg) //将image信息存入文件中
}
