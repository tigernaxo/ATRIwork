package fileformat

import (
	"testing"
)

func TestReadSingleFasta(t *testing.T) {
	_, seq := ReadSingleFasta("./data/test.fasta")
	t.Log(seq)
}

func TestSNVMapAllToAll(t *testing.T) {
	m, a := SNVMapAllToAll([][]byte{
		[]byte{65, 84, 67, 71, 97, 116, 67, 103, 45, 71},
		[]byte{65, 71, 67, 71, 97, 116, 67, 103, 45},
		[]byte{65, 84, 67, 78, 97, 116, 67, 103, 45},
		[]byte{65, 84, 67, 71, 97, 116, 67, 103, 110},
	})
	t.Log(m)
	t.Log(a)
}
func TestSNVsitesFromRef(t *testing.T) {
	ref := []byte("AtCcATCGATnGNTGCa")
	seq := []byte("ATCGATCGATCG-TGC")
	b := SNVsitesFromRef(ref, seq)
	t.Log(b)
	t.Logf("length of snvMap: %d\n", len(b))
}
