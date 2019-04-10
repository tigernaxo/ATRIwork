package fileformat

import (
	"bufio"
	"log"
	"os"
)

// ReadSingleFasta receive file name and return fasta content
func ReadSingleFasta(fasta string) (id string, sequence []byte) {

	f, err := os.Open(fasta)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	idByte := make([]byte, 0, 50)
	seq := make([]byte, 0, 10000)
	scanner := bufio.NewScanner(f)
	// > 62, \n10
	// - 45, A-Z 65-9, a-z 97-122
	isIDSec, isSeqSec := false, false

	for scanner.Scan() {
		for i, b := range scanner.Bytes() {
			// 每一行開頭需判斷目前讀取的是id或sequence
			if i == 0 {
				switch b {
				case 62:
					// 從頭進入 id section的情況
					if !isSeqSec && !isIDSec {
						isIDSec, isSeqSec = true, false
					}
				default:
					// 從id section 進入 sequence section的情況
					if isIDSec {
						isIDSec, isSeqSec = false, true
					}
				}
			}
			if isSeqSec {
				// 如果正在讀取sequence
				if IsAlphabet(b) || IsMisAlign(b) {
					seq = append(seq, b)
				} else {
					log.Panic("Error: Sequence must be alphabet or -")
				}
			} else if isIDSec {
				// 如果正在讀取id
				if i != 0 {
					idByte = append(idByte, b)
				}
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
