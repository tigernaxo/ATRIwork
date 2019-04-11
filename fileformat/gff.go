package fileformat

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Feature 保存特徵範圍、正負股、名稱
type Feature struct {
	Start  int    // gff col 4
	End    int    // gff col 5
	Name   string // gff col9 extract Name=name
	Strand byte   // gff col 7
}

// FeatureSet 保存同一類特徵、列表
type FeatureSet struct {
	Class    string // gff col 3
	Features []*Feature
}

// FeatureSetFromGFF 從GFF3裡面抽取特徵的範圍(start, end)、正負股、名稱
func FeatureSetFromGFF(featureClass, gffFile string) *FeatureSet {
	fs := FeatureSet{
		Class:    featureClass,
		Features: make([]*Feature, 0, 100000),
	}
	f, err := os.Open(gffFile)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	pattern := regexp.MustCompile(`Name=[^;$]*`)
	for scanner.Scan() {
		lineArr := strings.Split(scanner.Text(), "\t")
		if !strings.HasPrefix(lineArr[0], "#") && len(lineArr) == 9 && lineArr[2] == featureClass {
			start, err := strconv.Atoi(lineArr[3])
			if err != nil {
				log.Panic(err)
			}
			end, err := strconv.Atoi(lineArr[4])
			if err != nil {
				log.Panic(err)
			}
			fs.Features = append(fs.Features, &Feature{
				Start:  start,
				End:    end,
				Name:   strings.TrimPrefix(string(pattern.Find([]byte(lineArr[8]))), "Name="),
				Strand: byte(lineArr[6][0]),
			})
		}
	}
	return &fs
}
