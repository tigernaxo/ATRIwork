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
	fastas := os.Args[2:]
	fmt.Printf("Log: Reference GFF: %s\n", gff)
	fmt.Printf("Log: Sequence number: %d\n", len(fastas))

	seqs := make([][]byte, 0, len(fastas))

	c := make(chan []byte, len(fastas))
	for _, fa := range fastas {
		go func(f string, c chan []byte) {
			fmt.Printf("Log: Reading file: %s\n", f)
			_, seq := fileformat.ReadSingleFasta(f)
			// seqs = append(seqs, seq)
			c <- seq
		}(fa, c)
		// fmt.Printf("Log: Reading file: %s\n", fa)
		// id, seq := fileformat.ReadSingleFasta(fa)
		// seqs = append(seqs, seq)
		faByte := <-c
		seqs = append(seqs, faByte)

	}
	fmt.Printf("Log: Calculating SNV amoung all fasta...\n")
	siteMap, _ := snv.SiteMapAllToAll(seqs)

	fmt.Printf("Log: Extracting gene from GFF...\n")
	a := snv.SiteAnnotator{
		Sites:      snv.SiteMapToSiteSlice(siteMap),
		FeatureSet: fileformat.FeatureSetFromGFF("gene", gff),
	}

	fmt.Printf("Log: Annotating SNV ...\n")
	a.AnnotateAndSave("annotatedSNV.tsv")
}
