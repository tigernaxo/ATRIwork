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

// AnnotateAndSave 註解之後直接存檔
func (v *SiteAnnotator) AnnotateAndSave(fileName string) {

	file, err := os.Create(fileName)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	for _, feature := range v.FeatureSet.Features {
		for i := feature.Start - 1; i < feature.End; i++ {
			if v.SiteMap[i] {
				s := fmt.Sprintf("%d\t%s\t%s\t%d\t%d\n", i+1, feature.Name, string(feature.Strand), feature.Start, feature.End)
				_, err := file.WriteString(s)
				if err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
