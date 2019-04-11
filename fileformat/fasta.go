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
	seq := make([]byte, 0, 10000)
	// > 62, \n 10, - 45
	// A-Z 65-90, a-z 97-122
	isLineStart, isIDSec, isSeqSec := true, false, false
	// 行末\n              isLineStart = true
	// id section行末的\n  IDsec = false
	// seq section行末\n   isLineStart = true
	// 每行第一個byte是\n       isLineStart = true
	// 每行第一個byte是'>'       isIDsec = true
	// 每行第一個byte但不是>\n 如果是其他: 開始或繼續seq section
	// 每行第一個byte但不是>\n 如果是>: 結束seq section開始id section

	for _, b := range content {
		switch b {
		case 10:
			if isLineStart {
				// 空白行(^$)的情況
				isLineStart = true
			} else if isIDSec {
				// 到達ID section尾端
				isLineStart, isIDSec, isSeqSec = true, false, true
			} else if isSeqSec {
				// 在Seq section內換行
				isLineStart, isIDSec = true, false
			}
		case 62:
			if isLineStart {
				// >在第一個char就啟動id section
				isLineStart, isIDSec, isSeqSec = false, true, false
			} else {
				// >在其他位置就直接加上去
				seq = append(seq, b)
			}
		default:
			if isSeqSec && (IsAlphabet(b) || IsMisAlign(b)) {
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
