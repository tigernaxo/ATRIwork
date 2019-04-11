package snv

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("usage: SNPAnnotator VCF_FILE GFF_FILE OUTPUT_FILE\n")
		os.Exit(0)
	}
	// vcf := os.Args[1]
	// gff := os.Args[2]
	// out := os.Args[3]

	// v := SiteAnnotator{
	// 	Sites:      fileformat.SiteFromVCF(vcf),
	// 	FeatureSet: fileformat.FeatureSetFromGFF("gene", gff),
	// }
	// v.AnnotateAndSave(out)
}
