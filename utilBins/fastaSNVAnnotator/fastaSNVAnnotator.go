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

	seqs := make([][]byte, 0, len(fastas))
	ids := make([]string, 0, len(fastas))

	for _, fa := range fastas {
		id, seq := fileformat.ReadSingleFasta(fa)
		fmt.Printf("length of %s : %d\n", id, len(seq))
		seqs = append(seqs, seq)
		ids = append(ids, id)
	}
	siteMap, _ := snv.SiteMapAllToAll(seqs)
	a := snv.SiteAnnotator{
		Sites:      snv.SiteMapToSiteSlice(siteMap),
		FeatureSet: fileformat.FeatureSetFromGFF("gene", gff),
	}
	// fmt.Println(len(a.FeatureSet.Features))
	// fmt.Println(len(seqs))
	// fmt.Println(len(siteMap))
	// fmt.Println(len(a.Sites))
	a.AnnotateAndSave("annotatedSNV.tsv")
}
