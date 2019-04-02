package main

import (
	"testing"

	"github.com/tigernaxo/ATRIwork/fileformat"
)

func TestAnnotateAndSave(t *testing.T) {
	v := VCFAnnotator{
		SNPsites:   fileformat.SiteFromVCF("rawData/snp.vcf"),
		FeatureSet: fileformat.FeatureSetFromGFF("gene", "rawData/AE006468.gff3"),
	}
	v.AnnotateAndSave("out.tsv")
}
