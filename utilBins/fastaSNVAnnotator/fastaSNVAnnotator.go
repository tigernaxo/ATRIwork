package main

import (
	"fmt"
	"os"
	"time"

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

	for _, fa := range fastas {
		t := time.Now()
		fmt.Printf("[%02v:%02v:%02v] ", t.Hour(), t.Minute(), t.Second())
		fmt.Printf("Reading %s\n", fa)
		_, seq := fileformat.ReadSingleFasta(fa)
		siteMap = snv.SiteMapUpdate(siteMap, ref, seq)
	}

	t := time.Now()
	fmt.Printf("[%02v:%02v:%02v] ", t.Hour(), t.Minute(), t.Second())
	fmt.Printf("Extracting gene from %s...\n", gff)
	a := snv.SiteAnnotator{
		SiteMap:    siteMap,
		FeatureSet: fileformat.FeatureSetFromGFF("gene", gff),
	}

	t = time.Now()
	fmt.Printf("[%02v:%02v:%02v] ", t.Hour(), t.Minute(), t.Second())
	fmt.Printf("Annotating SNV ...\n")
	a.AnnotateAndSave("annotatedSNV.tsv")
}
