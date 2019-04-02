package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/tigernaxo/ATRIwork/fileformat"
)

// VCFAnnotator 的除非是使用bed等等0base的系統，定位皆為1base。
type VCFAnnotator struct {
	SNPsites   []int
	FeatureSet *fileformat.FeatureSet
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

// type Feature struct {
// 	Start  int    // gff col 4
// 	End    int    // gff col 5
// 	Name   string // gff col9 extract Name=name
// 	Strand byte   // gff col 7
// }

// // FeatureSet 保存同一類特徵、列表
// type FeatureSet struct {
// 	Class    string // gff col 3
// 	Features []*Feature
// }

// // SiteFromVCF 抽取VCF裡面的位址，回傳一個int slice
// func SiteFromVCF(vcf string) []int {
// 	fmt.Printf("Extracting site information from vcf file: [%s].\n", vcf)
// 	f, err := os.Open(vcf)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	defer f.Close()

// 	sites := make([]int, 0, 10000)
// 	scanner := bufio.NewScanner(f)
// 	// lineArr := make([]string, 2)
// 	var line string
// 	for scanner.Scan() {
// 		// do somethign with line (scanner.Text)
// 		line = scanner.Text()
// 		lineArr := strings.Split(line, "\t")
// 		if !strings.HasPrefix(lineArr[0], "#") {
// 			i, err := strconv.Atoi(lineArr[1])
// 			sites = append(sites, i)
// 			if err != nil {
// 				log.Panic(err)
// 			}
// 		}
// 	}

// 	return sites
// }

// // FeatureSetFromGFF 從GFF3裡面抽取特徵的範圍(start, end)、正負股、名稱
// func FeatureSetFromGFF(featureClass, gffFile string) *FeatureSet {
// 	fmt.Printf("Extracting feature class [%s] from file [%s]\n", featureClass, gffFile)
// 	fs := FeatureSet{
// 		Class:    featureClass,
// 		Features: make([]*Feature, 0, 100000),
// 	}
// 	f, err := os.Open(gffFile)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	defer f.Close()

// 	scanner := bufio.NewScanner(f)
// 	pattern := regexp.MustCompile(`Name=[^;$]*`)
// 	for scanner.Scan() {
// 		lineArr := strings.Split(scanner.Text(), "\t")
// 		if !strings.HasPrefix(lineArr[0], "#") && len(lineArr) == 9 && lineArr[2] == featureClass {
// 			start, err := strconv.Atoi(lineArr[3])
// 			if err != nil {
// 				log.Panic(err)
// 			}
// 			end, err := strconv.Atoi(lineArr[4])
// 			if err != nil {
// 				log.Panic(err)
// 			}
// 			fs.Features = append(fs.Features, &Feature{
// 				Start:  start,
// 				End:    end,
// 				Name:   strings.TrimPrefix(string(pattern.Find([]byte(lineArr[8]))), "Name="),
// 				Strand: byte(lineArr[6][0]),
// 			})
// 		}
// 	}
// 	return &fs
// }
