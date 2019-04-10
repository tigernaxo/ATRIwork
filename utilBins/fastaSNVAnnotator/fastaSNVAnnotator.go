package main

import (
	"fmt"
	"os"

	"github.com/tigernaxo/ATRIwork/fileformat"
	// github.com/tigernaxo/ATRIwork/histogram
	"github.com/tigernaxo/ATRIwork/snv"
)

// Input:
// 1. intensity
// 2. ref gff
// 3. ...seqs

// Fixed:
// 4.color
// 5.bgcolor
// 6.min Unit ( ex:1000000, to decide convasLength )
// 7.out dimension
// 8.out name

// Output
// 1.histograms

// local import need
// github.com/tigernaxo/ATRIwork/fileformat
// github.com/tigernaxo/ATRIwork/histogram
// github.com/tigernaxo/ATRIwork/snv

func main() {
	gff := os.Args[1]
	fmt.Printf("Log: Reading file: %s\n", os.Args[2])
	_, ref := fileformat.ReadSingleFasta(os.Args[2])
	fastas := os.Args[3:]
	fmt.Printf("Log: Reference GFF: %s\n", gff)
	fmt.Printf("Log: Sequence number: %d\n", len(fastas))

	siteMap := make([]bool, len(ref))
	for i := range siteMap {
		siteMap[i] = false
	}
	// TODO:不應該用[][]byte
	// 應該seq比對完就丟棄，效能才會好
	// 但是這樣就不能用SiteMapAllToAll
	// 應該寫一個SiteMapUpdate(siteMap *[]bool, ref *[]byte, seq *[]byte)(newSiteMap []bool)
	// library裡面都是傳值呼叫，也要修改
	for _, fa := range fastas {
		fmt.Printf("Log: Reading file: %s\n", fa)
		_, seq := fileformat.ReadSingleFasta(fa)
		siteMap = snv.SiteMapUpdate(siteMap, ref, seq)
	}

	fmt.Printf("Log: Extracting gene from GFF...\n")
	a := snv.SiteAnnotator{
		Sites:      snv.SiteMapToSiteSlice(siteMap),
		FeatureSet: fileformat.FeatureSetFromGFF("gene", gff),
	}

	fmt.Printf("Log: Annotating SNV ...\n")
	a.AnnotateAndSave("annotatedSNV.tsv")
}
