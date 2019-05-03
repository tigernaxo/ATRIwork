package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tigernaxo/ATRIwork/snv"

	"github.com/tigernaxo/ATRIwork/fileformat"
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
	// Using ref length as map length
	_, r := fileformat.ReadSingleFasta(os.Args[2])
	mapLen := len(r)
	r = nil

	// Using first seq as ref (not reference fasta)
	// _, ref := fileformat.ReadSingleFasta(os.Args[3])

	// set sequence fasta
	fastas := os.Args[3:]
	fmt.Printf("Sequence number: %d\n", len(fastas))

	// Calculate each snv siteinfo
	maskChar := []byte{'N', '-'}
	ntToCount := []byte{'A', 'T', 'C', 'G'}
	ntCounter := make(map[byte][]uint8)
	for _, nt := range ntToCount {
		ntCounter[nt] = make([]uint8, mapLen)
		for j := range ntCounter[nt] {
			ntCounter[nt][j] = 0
		}
	}

	snvMap := make([]bool, mapLen)
	showMap := make([]bool, mapLen)
	finalMap := make([]bool, mapLen)
	fillBoolArr(showMap, true)
	fillBoolArr(snvMap, false)
	var seq, ref []byte
	for i, fa := range fastas {
		fmt.Printf("%s Reading %s\n", timeStamp(), fa)
		_, seq = fileformat.ReadSingleFasta(fa)
		if i == 0 {
			_, ref = fileformat.ReadSingleFasta(fa)
		}
		snv.UpdateSNVMapShowMapSiteCount(ref, seq, maskChar, snvMap, showMap, ntCounter)
	}
	// Debug
	// fmt.Println("==============Debug Start==============")
	// for i := 0; i < len(ntCounter['A']); i++ {
	// 	sum := 0
	// 	site := i + 1
	// 	for _, arr := range ntCounter {
	// 		sum = sum + int(arr[i])
	// 	}
	// 	if sum != len(fastas) && showMap[i] {
	// 		fmt.Printf("%d\n", i)
	// 		fmt.Printf("Length of counter array: %d \n", len(ntCounter['A']))
	// 		fmt.Printf("Site number:\n")
	// 		fmt.Printf("\tA:%d", ntCounter['A'][site])
	// 		fmt.Printf("\tT:%d", ntCounter['T'][site])
	// 		fmt.Printf("\tC:%d", ntCounter['C'][site])
	// 		fmt.Printf("\tG:%d", ntCounter['G'][site])
	// 		fmt.Printf("\n")
	// 		fmt.Printf("Site: %d(1base) bases is:\n", i+1)
	// 		for _, fa := range fastas {
	// 			_, seq := fileformat.ReadSingleFasta(fa)
	// 			fmt.Printf("\t%c", seq[site])
	// 		}
	// 		fmt.Printf("\n")
	// 		break
	// 	}
	// }
	// fmt.Println("==============Debug End==============")
	// Debug

	// Decide final show site
	snv.BoolArrAND(finalMap, showMap, snvMap)

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
			if finalMap[i] {
				s := fmt.Sprintf("%d\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%c\n",
					i+1, feature.Name, string(feature.Strand), feature.Start, feature.End,
					ntCounter['A'][i], ntCounter['T'][i], ntCounter['C'][i], ntCounter['G'][i], ref[i])
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

	aln := make([]byte, 0, logBoolCount(finalMap, true))
	for _, fa := range fastas {
		id, seq := fileformat.ReadSingleFasta(fa)
		// write id
		_, err := f.WriteString(">" + id + "\n")
		logErr(err)

		// write single alignment
		for i, c := range seq {
			if finalMap[i] {
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

func fillBoolArr(arr []bool, fill bool) {
	for i := range arr {
		arr[i] = fill
	}
}
