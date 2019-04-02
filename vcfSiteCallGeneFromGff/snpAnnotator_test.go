package main

import (
	"testing"
)

func TestSiteFromVCF(t *testing.T) {
	siteCount := len(SiteFromVCF("rawData/snp.vcf"))
	if siteCount != 129975 {
		t.Error("vcf sites count:", siteCount)
		t.Fail()
	}
	t.Log("vcf sites count:", siteCount)
}
func TestFeatureSetFromGFF(t *testing.T) {
	fs := FeatureSetFromGFF("gene", "rawData/AE006468.gff3")
	t.Error("fs class:", fs.Class)
}
func TestAnnotateAndSave(t *testing.T) {
	v := VCFAnnotator{
		SNPsites:   SiteFromVCF("rawData/snp.vcf"),
		FeatureSet: FeatureSetFromGFF("gene", "rawData/AE006468.gff3"),
	}
	v.AnnotateAndSave("out.tsv")
}
