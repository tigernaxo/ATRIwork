package main

import (
	"testing"
)

func TestAnnotateAndSave(t *testing.T) {
	v := VCFAnnotator{
		SNPsites:   SiteFromVCF("rawData/snp.vcf"),
		FeatureSet: FeatureSetFromGFF("gene", "rawData/AE006468.gff3"),
	}
	v.AnnotateAndSave("out.tsv")
}
