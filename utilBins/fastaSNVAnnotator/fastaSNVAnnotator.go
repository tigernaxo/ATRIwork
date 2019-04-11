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
	fmt.Printf("Reading file: %s\n", os.Args[2])
	_, ref := fileformat.ReadSingleFasta(os.Args[2])
	fastas := os.Args[3:]
	fmt.Printf("Reference GFF: %s\n", gff)
	fmt.Printf("Sequence number: %d\n", len(fastas))

	siteMap := make([]bool, len(ref))
	for i := range siteMap {
		siteMap[i] = false
	}

	for _, fa := range fastas {
		fmt.Printf("%s Reading %s\n", timeStamp(), fa)
		_, seq := fileformat.ReadSingleFasta(fa)
		siteMap = snv.SiteMapUpdate(siteMap, ref, seq)
	}

	fmt.Printf("%s Extracting gene from %s...\n", timeStamp(), gff)
	a := snv.SiteAnnotator{
		SiteMap:    siteMap,
		FeatureSet: fileformat.FeatureSetFromGFF("gene", gff),
	}

	fmt.Printf("%s Annotating SNV ...\n", timeStamp())
	a.AnnotateAndSave("annotatedSNV.tsv")

	fmt.Printf("%s Finished!\n", timeStamp())
}

func timeStamp() string {
	t := time.Now()
	return fmt.Sprintf("[%02v:%02v:%02v]", t.Hour(), t.Minute(), t.Second())
}
