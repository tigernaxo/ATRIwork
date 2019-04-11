package snv

import (
	"fmt"
	"log"
	"os"

	"github.com/tigernaxo/ATRIwork/fileformat"
)

// SiteAnnotator 的除非是使用bed等等0base的系統，定位皆為1base。
type SiteAnnotator struct {
	SiteMap    []bool
	FeatureSet *fileformat.FeatureSet
}

// SiteMapToSiteSlice 把SiteMap轉成Sites
func SiteMapToSiteSlice(siteMap []bool) (siteSlice []int) {
	siteSlice = make([]int, 0, len(siteMap))
	for i, b := range siteMap {
		if b {
			siteSlice = append(siteSlice, i+1)
		}
	}
	return siteSlice
}

// AnnotateAndSave 註解之後直接存檔
// func (v *SiteAnnotator) AnnotateAndSave(fileName string) {

// 	fmt.Println("Sorting snp sites...")
// 	sort.Ints(v.Sites)

// 	fmt.Printf("Creating file: [%s]\n", fileName)
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	defer file.Close()

// 	for i := range v.Sites {
// 		for _, f := range v.FeatureSet.Features {
// 			if f.Start <= i && f.End >= i {
// 				s := fmt.Sprintf("%d\t%s\t%s\t%d\t%d\n", i, f.Name, string(f.Strand), f.Start, f.End)
// 				_, err := file.WriteString(s)
// 				if err != nil {
// 					log.Panic(err)
// 				}
// 			}
// 		}
// 	}
// }

// AnnotateAndSave 註解之後直接存檔
func (v *SiteAnnotator) AnnotateAndSave(fileName string) {

	fmt.Printf("Creating file: [%s]\n", fileName)
	file, err := os.Create(fileName)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	for _, feature := range v.FeatureSet.Features {
		for i := feature.Start - 1; i < feature.End; i++ {
			if v.SiteMap[i] {
				s := fmt.Sprintf("%d\t%s\t%s\t%d\t%d\n", i, feature.Name, string(feature.Strand), feature.Start, feature.End)
				_, err := file.WriteString(s)
				if err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
