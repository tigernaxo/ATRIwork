package main

import (
	"fmt"
	"log"
	"os"
	"sort"
)

// VCFAnnotator 的除非是使用bed等等0base的系統，定位皆為1base。
type VCFAnnotator struct {
	SNPsites   []int
	FeatureSet *FeatureSet
}

// AnnotateAndSave 註解之後直接存檔
func (v *VCFAnnotator) AnnotateAndSave(fileName string) {

	fmt.Println("Sorting snp sites...")
	sort.Ints(v.SNPsites)

	fmt.Printf("Creating file: [%s]\n", fileName)
	file, err := os.Create(fileName)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	for i := range v.SNPsites {
		for _, f := range v.FeatureSet.Features {
			if f.Start <= i && f.End >= i {
				s := fmt.Sprintf("%d\t%s\t%s\t%d\t%d\n", i, f.Name, string(f.Strand), f.Start, f.End)
				_, err := file.WriteString(s)
				if err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
