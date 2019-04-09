package fileformat

import (
	"testing"
)

func TestReadSingleFasta(t *testing.T) {
	_, seq := ReadSingleFasta("./data/test.fasta")
	t.Error(seq)
}
