package snv

import (
	"testing"

	"github.com/tigernaxo/ATRIwork/fileformat"
)

func TestAnnotateAndSave(t *testing.T) {
	v := SiteAnnotator{
		sites:      fileformat.SiteFromVCF("data/snp.vcf"),
		FeatureSet: fileformat.FeatureSetFromGFF("gene", "data/AE006468.gff3"),
	}
	v.AnnotateAndSave("out.tsv")
}
