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

// SNVMapAllToAll take slice of id, seq byte and return snv sites, snv alignment
// SNVMapAllToAll treat - different from n/N, so - and n on the same site are consider as snv.
func SNVMapAllToAll(seqs [][]byte) (snvMap []bool, snvAlign [][]byte) {
	// 先判斷最大長度是多少
	var maxLength int
	for _, seq := range seqs {
		if maxLength < len(seq) {
			maxLength = len(seq)
		}
	}

	snvAlign = make([][]byte, 0, len(seqs))
	snvMap = make([]bool, maxLength)
	for i := range snvMap {
		snvMap[i] = false
	}

	// 以第一個序列作為參考序列，並填充-到末端
	ref := make([]byte, maxLength)
	if refLength := copy(ref, seqs[0]); refLength < maxLength {
		for i := refLength; i < maxLength; i++ {
			ref[i] = 45
			snvMap[i] = true
		}
	}

	// 逐序列比對
	for _, seq := range seqs {

		// 逐nt比對

		// 處理len(seq)範圍內的部份
		for i, nt := range seq {
			// 取得snv map
			if !snvMap[i] && (nt-ref[i])%32 != 0 {
				snvMap[i] = true
			}
		}

		// 處理大於seq小於maxLength的部份
		for tailSpaceIndex := len(seq); tailSpaceIndex < maxLength; tailSpaceIndex++ {
			// 如果ref不是-(ref有序列而seq沒有)就判斷該位置是snv
			if !snvMap[tailSpaceIndex] && ref[tailSpaceIndex] != 45 {
				snvMap[tailSpaceIndex] = true
			}
		}
	}
	// 取得snvMap內true的數量
	SNVnumber := 0
	for _, b := range snvMap {
		if b {
			SNVnumber++
		}
	}

	// 根據snv map取得alignment
	// 並且push到snvAlign
	for _, seq := range seqs {
		seqSNV := make([]byte, 0, SNVnumber)
		for i := range seq {
			if snvMap[i] {
				seqSNV = append(seqSNV, seq[i])
			}
		}
		for seqLength := len(seq); seqLength < maxLength; seqLength++ {
			seqSNV = append(seqSNV, 45)
		}
		snvAlign = append(snvAlign, seqSNV)
	}

	return snvMap, snvAlign
}

// SNVsitesFromRef to be finish
func SNVsitesFromRef(ref []byte, seq []byte) (snvMap []bool) {

	// 如果seq的長度比ref長就跳錯
	if len(ref) < len(seq) {
		log.Panic("Error: target sequence length is longer then reference")
	}

	// 先做一個空的snvMap並填滿false
	snvMap = make([]bool, len(ref))
	for i := range snvMap {
		snvMap[i] = false
	}

	// 找出snv
	for i, nt := range seq {
		if (nt-ref[i])%32 != 0 {
			snvMap[i] = true
		}
	}

	// 處理reference多出來的部分
	for i := len(seq); i < len(ref); i++ {
		if !IsMisAlign(ref[i]) {
			snvMap[i] = true
		}
	}

	// 回傳snvMap
	return snvMap
}

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
