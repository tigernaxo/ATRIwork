package snv

import "testing"

func TestSiteMapAllToAll(t *testing.T) {
	m, a := SiteMapAllToAll([][]byte{
		[]byte{65, 84, 67, 71, 97, 116, 67, 103, 45, 71},
		[]byte{65, 71, 67, 71, 97, 116, 67, 103, 45},
		[]byte{65, 84, 67, 78, 97, 116, 67, 103, 45},
		[]byte{65, 84, 67, 71, 97, 116, 67, 103, 110},
	})
	t.Log(m)
	t.Log(a)
}
func TestSiteMapFromRef(t *testing.T) {
	ref := []byte("AtCcATCGATnGNTGCa")
	seq := []byte("ATCGATCGATCG-TGC")
	b := SiteMapFromRef(ref, seq)
	t.Log(b)
	t.Logf("length of snvMap: %d\n", len(b))
}
