package snv

import (
	"log"

	"github.com/tigernaxo/ATRIwork/fileformat"
)

// SiteInfo struct contain accumulate number of Nt and ref Nt
type SiteInfo struct {
	RefSeq                        []byte
	seqCount                      uint
	firstSeq                      []byte
	NtCount                       map[byte][]uint8
	SnvMap                        []bool
	ShowMask                      []bool
	ExcludeSiteIfAnySampleContain []byte
}

// NewSiteInfo return a new SiteInfo
func NewSiteInfo(ref []byte, ntSet []byte, excludeSideByNtSet []byte) *SiteInfo {
	// Default ShowMask (all true)
	mask := make([]bool, len(ref))
	for i := range mask {
		mask[i] = true
	}
	// Default SnvMap (all false)
	snvMap := make([]bool, len(ref))
	for i := range snvMap {
		snvMap[i] = false
	}
	// Default NtCount
	ntCount := make(map[byte][]uint8, len(ntSet))
	for _, nt := range ntSet {
		tmp := make([]uint8, len(ref))
		for i := range tmp {
			tmp[i] = 0
		}
		ntCount[nt] = tmp
	}

	// return SiteInfo
	return &SiteInfo{
		seqCount:                      0,
		RefSeq:                        ref,
		ShowMask:                      mask,
		SnvMap:                        snvMap,
		NtCount:                       ntCount,
		ExcludeSiteIfAnySampleContain: excludeSideByNtSet,
	}
}

// AccumulateSNV accumulate sequence to SiteInfo
func (siteInfo *SiteInfo) AccumulateSNV(seq []byte) {
	siteInfo.seqCount++
	if siteInfo.seqCount == 1 {
		siteInfo.firstSeq = seq
	}
	// Panic if seq length > reference length
	if len(siteInfo.RefSeq) < len(seq) {
		log.Panicf("[Error] Sequence No.%d length %d is longer than referrnce sequence %d\n", siteInfo.seqCount, len(seq), len(siteInfo.RefSeq))
	}
	for i, c := range seq {
		// Update SnvMap
		if siteInfo.seqCount != 1 {
			if !IsEqualAlphabet(c, siteInfo.firstSeq[i]) {
				siteInfo.SnvMap[i] = true
			}
		}
		// Update NtCount
		for nt, nta := range siteInfo.NtCount {
			if IsEqualAlphabet(c, nt) {
				nta[i]++
			}
		}
		// Update ShowMask
		for _, ex := range siteInfo.ExcludeSiteIfAnySampleContain {
			if IsEqualAlphabet(c, ex) {
				siteInfo.ShowMask[i] = false
			}
		}
	}
}

// SnvmapAndMask output And operated slice
func (siteInfo *SiteInfo) SnvmapAndMask(outMap []bool) {
	if len(outMap) != len(siteInfo.SnvMap) {
		log.Panicf("[Error] Trying to calculate (SnvMap AND Mask), but ")
		log.Panicf("the length of output []boole vector is not equal to reference length\n")
	}
	for i, b := range siteInfo.SnvMap {
		if siteInfo.ShowMask[i] && b {
			outMap[i] = true
		} else {
			outMap[i] = false
		}
	}
}

func IsAlphabet(c byte) bool {
	if 65 <= c && c <= 122 {
		if c <= 90 || c >= 97 {
			return true
		}
	}
	return false
}
func IsEqualAlphabet(c, d byte) bool {
	if (c-d)%32 == 0 {
		if IsAlphabet(c) && IsAlphabet(d) {
			return true
		}
	}
	return false
}

// ==============old=========================

// UpdateCharMap update charMap, set true when site's char belong charSet
func UpdateCharMap(charMap []bool, seq []byte, charSet []byte) {
	// If len(siteMap) < len(seq) , append siteMap and set append section all false
	if len(charMap) < len(seq) {
		appendSlice := make([]bool, len(seq)-len(charMap))
		for i := range appendSlice {
			appendSlice[i] = false
		}
		charMap = append(charMap, appendSlice...)
	}
	for i, c := range seq {
		for _, d := range charSet {
			if c == d {
				charMap[i] = true
				continue
			}
		}
	}
}

// UpdateSNVMap 根據ref和seq更新siteMap
func UpdateSNVMap(siteMap []bool, ref []byte, seq []byte) {
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
