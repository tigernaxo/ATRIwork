package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tigernaxo/ATRIwork/fileformat"
	"github.com/tigernaxo/ATRIwork/snv"
)

// Input:
// 1. ref gff
// 2. ...seqs

// Output
// 1.annotation tsv

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Printf("\t%s reference.gff ref.fa seq1.fa seq2.fa [seq3.fa ...]\n", os.Args[0])
		os.Exit(0)
	}
	gff := os.Args[1]
	fmt.Printf("Reference GFF: %s\n", gff)
	fmt.Printf("Reading file as reference: %s\n", os.Args[2])
	_, ref := fileformat.ReadSingleFasta(os.Args[2])

	fastas := os.Args[3:]
	fmt.Printf("Sequence number: %d\n", len(fastas))

	siteInfo := snv.NewSiteInfo(ref, []byte{'A', 'T', 'C', 'G'}, []byte{'N'})
	var seq []byte
	for _, fa := range fastas {
		fmt.Printf("%s Reading %s\n", timeStamp(), fa)
		_, seq = fileformat.ReadSingleFasta(fa)
		siteInfo.AccumulateSNV(seq)
	}
	showMap := make([]bool, len(ref))
	siteInfo.SnvmapAndMask(showMap)

	fmt.Printf("%s Extracting gene from %s...\n", timeStamp(), gff)
	a := snv.SiteAnnotator{
		SiteMap:    showMap,
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

func logBoolCount(name string, ba []bool, checkB bool) {
	count := 0
	for _, b := range ba {
		if b == checkB {
			count++
		}
	}
	fmt.Printf("%s %v count : %d\n", name, checkB, count)
}
