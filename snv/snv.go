package snv

import (
	"log"

	"github.com/tigernaxo/ATRIwork/fileformat"
)

// SiteMapUpdate 根據ref和seq更新siteMap
func SiteMapUpdate(siteMap []bool, ref []byte, seq []byte) (newSiteMap []bool) {
	// 如果seq > siteMap就要append
	// 要注意ref out of index
	if len(seq) > len(siteMap) {
		newMap := make([]bool, 0, len(seq))
		for i := copy(newMap, siteMap); i < len(seq); i++ {
			switch seq[i] {
			case 45:
				newMap[i] = false
			default:
				newMap[i] = true
			}
		}
		siteMap = newMap
	}

	// 如果數字不同就true siteMap
	for i := range seq {
		if !siteMap[i] && (ref[i]-seq[i])%32 != 0 {
			siteMap[i] = true
		}
	}
	return siteMap
}

// SiteMapAllToAll take slice of id, seq byte and return snv sites, snv alignment
// SiteMapAllToAll treat - different from n/N, so - and n on the same site are consider as snv.
func SiteMapAllToAll(seqs [][]byte) (snvMap []bool, snvAlign [][]byte) {
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

// SiteMapFromRef to be finish
func SiteMapFromRef(ref []byte, seq []byte) (snvMap []bool) {

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
		if !fileformat.IsMisAlign(ref[i]) {
			snvMap[i] = true
		}
	}

	// 回傳snvMap
	return snvMap
}
