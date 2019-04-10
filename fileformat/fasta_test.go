package fileformat

import (
	"testing"
)

func TestReadSingleFasta(t *testing.T) {
	id, seq := ReadSingleFasta("../testData/S17-034-AE006468.2.SNV.fasta")
	t.Error(id)
	t.Error(len(seq))
}
