package fileformat

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// SiteFromVCF 抽取VCF裡面的位址，回傳一個int slice
func SiteFromVCF(vcf string) []int {
	fmt.Printf("Extracting site information from vcf file: [%s].\n", vcf)
	f, err := os.Open(vcf)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	sites := make([]int, 0, 10000)
	scanner := bufio.NewScanner(f)
	// lineArr := make([]string, 2)
	var line string
	for scanner.Scan() {
		// do somethign with line (scanner.Text)
		line = scanner.Text()
		lineArr := strings.Split(line, "\t")
		if !strings.HasPrefix(lineArr[0], "#") {
			i, err := strconv.Atoi(lineArr[1])
			sites = append(sites, i)
			if err != nil {
				log.Panic(err)
			}
		}
	}

	return sites
}
