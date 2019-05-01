package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tigernaxo/ATRIwork/fileformat"
	"github.com/tigernaxo/ATRIwork/snv"
)

func main() {
	// output help message
	if len(os.Args) < 5 {
		fmt.Println("Usage:")
		fmt.Printf("\t%s reference.gff ref.fa seq1.fa seq2.fa [seq3.fa ...]\n", os.Args[0])
		os.Exit(0)
	}
	// set reference fasta, gff
	gff := os.Args[1]
	fmt.Printf("Reference GFF: %s\n", gff)
	fmt.Printf("Reading file as reference: %s\n", os.Args[2])
	_, ref := fileformat.ReadSingleFasta(os.Args[2])

	// set sequence fasta
	fastas := os.Args[3:]
	fmt.Printf("Sequence number: %d\n", len(fastas))

	// Calculate each snv siteinfo
	siteInfo := snv.NewInfo(ref, []byte{'N', '-'}, false)

	// Create nt counter
	ntToCount := []byte{'A', 'T', 'C', 'G'}
	ntCount := make(map[byte][]uint8)
	for _, nt := range ntToCount {
		ntCount[nt] = make([]uint8, len(ref))
		for i := range ntCount[nt] {
			ntCount[nt][i] = 0
		}
	}

	var seq []byte
	for _, fa := range fastas {
		fmt.Printf("%s Reading %s\n", timeStamp(), fa)
		_, seq = fileformat.ReadSingleFasta(fa)
		siteInfo.AccumulateSeqSNV(seq)
		// Count site nt frequency
		for i, nt := range seq {
			for _, c := range ntToCount {
				if snv.IsEqualAlphabet(c, nt) {
					ntCount[c][i]++
				}
			}
		}
	}

	// Decide show site
	showMap := make([]bool, len(ref))
	siteInfo.SnvmapAndMask(showMap)
	fmt.Printf("%s Totel SNV between sequences : %d\n", timeStamp(), logBoolCount(showMap, true))

	// extract gene feature from gff
	fmt.Printf("%s Extracting gene from %s...\n", timeStamp(), gff)
	var geneSet *fileformat.FeatureSet
	geneSet = fileformat.FeatureSetFromGFF("gene", gff)

	// Start output
	fmt.Printf("%s Creating annotated.tsv\n", timeStamp())
	file, err := os.Create("annotated.tsv")
	logErr(err)
	defer file.Close()

	// write column title
	titles := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		"Site", "FeatureName", "Strand", "Start", "End",
		"A_count", "T_count", "C_count", "G_count", "RefNt")
	_, err = file.WriteString(titles)
	logErr(err)

	// write feature set which range contain showMap true value
	for _, feature := range geneSet.Features {
		for i := feature.Start - 1; i < feature.End; i++ {
			if showMap[i] {
				s := fmt.Sprintf("%d\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%c\n",
					i+1, feature.Name, string(feature.Strand), feature.Start, feature.End,
					ntCount['A'][i], ntCount['T'][i], ntCount['C'][i], ntCount['G'][i], siteInfo.RefSeq[i])
				_, err := file.WriteString(s)
				logErr(err)
			}
		}
	}
	// Output alignment
	fmt.Printf("%s Outputing %s\n", timeStamp(), "alignment")
	f, err := os.Create("snv.aln")
	logErr(err)
	defer f.Close()

	aln := make([]byte, 0, logBoolCount(showMap, true))
	for _, fa := range fastas {
		id, seq := fileformat.ReadSingleFasta(fa)
		// write id
		_, err := f.WriteString(">" + id + "\n")
		logErr(err)

		// write single alignment
		for i, c := range seq {
			if showMap[i] {
				aln = append(aln, c)
			}
		}
		_, err = f.Write(aln)
		logErr(err)
		_, err = f.WriteString("\n")
		logErr(err)

		// empty alignment
		aln = aln[:0]
	}

	// Program finished
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
func logBoolCount(ba []bool, checkB bool) int {
	count := 0
	for _, b := range ba {
		if b == checkB {
			count++
		}
	}
	return count
}
