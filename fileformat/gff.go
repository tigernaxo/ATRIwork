package fileformat

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// FeatureSetFromGFF 從GFF3裡面抽取特徵的範圍(start, end)、正負股、名稱
func FeatureSetFromGFF(featureClass, gffFile string) *FeatureSet {
	fmt.Printf("Extracting feature class [%s] from file [%s]\n", featureClass, gffFile)
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
