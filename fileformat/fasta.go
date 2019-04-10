package fileformat

import (
	"io/ioutil"
	"log"
)

// ReadSingleFasta receive file name and return fasta content
func ReadSingleFasta(fasta string) (id string, sequence []byte) {

	content, err := ioutil.ReadFile(fasta)
	if err != nil {
		log.Fatal(err)
	}

	idByte := make([]byte, 0, 50)
	seq := make([]byte, 0, 10000) // > 62, \n10
	// - 45, A-Z 65-9, a-z 97-122
	isLineStart, isIDSec, isSeqSec := true, false, false

	for _, b := range content {
		// if i == 4857462 {
		// 	log.Panicf("i reaches %d, char is %d\n", i, int(b))
		// }
		switch b {
		case 10:
			if isLineStart {
				isLineStart, isIDSec, isSeqSec = true, false, false
			} else if isIDSec {
				isLineStart, isIDSec, isSeqSec = true, false, true
			} else {
				isLineStart, isIDSec = true, false
			}
		case 62:
			if isLineStart {
				isLineStart, isIDSec, isSeqSec = false, true, false
			} else {
				seq = append(seq, b)
			}
		default:
			if isSeqSec && (IsAlphabet(b) || IsMisAlign(b)) {
				// log.Panicf("i reaches %d, char is %d\n", i, int(b))
				seq = append(seq, b)
			} else if isIDSec {
				idByte = append(idByte, b)
			}
		}
	}
	return string(idByte), seq
}

// ReadMultiFasta receive file name and return fasta content
// func ReadMultiFasta(fasta string) (id []string, sequences [][]byte)

// IsSingleFasta receive file name and check weather the file is fasta or not
// func IsSingleFasta(f string) bool {
// 	return true
// }

// IsAlphabet take sequence byte and tell weather is a dna sequence or not
func IsAlphabet(b byte) bool {
	// - 45, A-Z 65-9, a-z 97-122
	if b >= 65 && b <= 122 {
		if b <= 90 || b >= 97 {
			return true
		}
	}
	return false
}

// IsMisAlign take sequence byte and tell weather is a mis align sequence or not
func IsMisAlign(b byte) bool {
	// - 45, A-Z 65-9, a-z 97-122
	if b == 45 {
		return true
	}
	return false
}
