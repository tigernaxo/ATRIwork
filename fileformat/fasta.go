package fileformat

import (
	"bufio"
	"log"
	"os"
)

// ReadFasta receive file name and return fasta content
func ReadFasta(fasta string) (id string, sequence []byte) {

	f, err := os.Open(fasta)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	idByte := make([]byte, 0, 50)
	seq := make([]byte, 0, 10000)
	scanner := bufio.NewScanner(f)
	// > 62, \n10
	isIDSec, isSeqSec := false, false

	for scanner.Scan() {
		for _, b := range scanner.Bytes() {
			if isSeqSec && (IsDNASequence(b) || IsMisAlign(b)) {
				seq = append(seq, b)
			} else if isIDSec && b == 10 {
				isIDSec, isSeqSec = false, true
				// case b...
			} else if b == 62 {
				isIDSec, isSeqSec = true, false
			}
		}
	}

	return "", seq
}

// IsSingleFasta receive file name and check weather the file is fasta or not
func IsSingleFasta(f string) bool {
	return true
}

// IsDNASequence take sequence byte and tell weather is a dna sequence or not
func IsDNASequence(b byte) bool {
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
