package main

import (
	"fmt"
	"log"
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
	if len(os.Args) < 5 {
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
	var geneSet *fileformat.FeatureSet
	geneSet = fileformat.FeatureSetFromGFF("gene", gff)

	file, err := os.Create("annotated.tsv")
	logErr(err)
	defer file.Close()

	fmt.Printf("%s Annotating SNV ...\n", timeStamp())

	// Column Title
	s := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		"Site", "FeatureName", "Strand", "Start", "End",
		"A_count", "T_count", "C_count", "G_count", "RefNt")
	_, err = file.WriteString(s)
	logErr(err)

	logErr(err)
	for _, feature := range geneSet.Features {
		for i := feature.Start - 1; i < feature.End; i++ {
			if showMap[i] {
				s := fmt.Sprintf("%d\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%c\n",
					i+1, feature.Name, string(feature.Strand), feature.Start, feature.End,
					siteInfo.NtCount['A'][i], siteInfo.NtCount['T'][i], siteInfo.NtCount['C'][i], siteInfo.NtCount['G'][i], siteInfo.RefSeq[i])
				_, err := file.WriteString(s)
				logErr(err)
			}
		}
	}
	fmt.Printf("%s All Finished!!\n", timeStamp())
}
func timeStamp() string {
	t := time.Now()
	return fmt.Sprintf("[%02v:%02v:%02v]", t.Hour(), t.Minute(), t.Second())
}
func logErr(e error) {
	if e != nil {
		log.Panic(e)
	}
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
