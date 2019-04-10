package snv

import (
	"testing"

	"github.com/tigernaxo/ATRIwork/fileformat"
)

func TestAnnotateAndSave(t *testing.T) {
	v := SiteAnnotator{
		Sites:      fileformat.SiteFromVCF("../testData/snp.vcf"),
		FeatureSet: fileformat.FeatureSetFromGFF("gene", "../testData/AE006468.gff3"),
	}
	v.AnnotateAndSave("out.tsv")
}
