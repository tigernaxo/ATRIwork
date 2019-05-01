package snv

import (
	"log"
)

// Info struct contain accumulate number of Nt and ref Nt
type Info struct {
	RefSeq                        []byte
	seqCount                      uint
	realRef                       []byte
	snvToRef                      bool
	SnvMap                        []bool
	ShowMask                      []bool
	ExcludeSiteIfAnySampleContain []byte
}

// NewInfo return a new SiteInfo
func NewInfo(ref []byte, excludeSideByNtSet []byte, snvToRef bool) *Info {
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

	// return SiteInfo
	return &Info{
		seqCount:                      0,
		RefSeq:                        ref,
		ShowMask:                      mask,
		SnvMap:                        snvMap,
		ExcludeSiteIfAnySampleContain: excludeSideByNtSet,
		snvToRef:                      snvToRef,
	}
}

// AccumulateSeqSNV accumulate sequence to SiteInfo
func (s *Info) AccumulateSeqSNV(seq []byte) {
	//
	if s.snvToRef {
		s.realRef = s.RefSeq
	} else {
		s.seqCount++
		if s.seqCount == 1 {
			s.realRef = seq
		}
	}

	// Panic if seq length > reference length
	if len(s.RefSeq) < len(seq) {
		log.Panicf("[Error] Sequence No.%d length %d is longer than referrnce sequence %d\n", s.seqCount, len(seq), len(s.RefSeq))
	}
	for i, c := range seq {
		// Update SnvMap
		if s.seqCount != 1 {
			if !IsEqualAlphabet(c, s.realRef[i]) {
				s.SnvMap[i] = true
			}
		}
		for _, ex := range s.ExcludeSiteIfAnySampleContain {
			if IsEqualAlphabet(c, ex) {
				s.ShowMask[i] = false
			}
		}
	}
}

// SnvmapAndMask output And operated slice
func (s *Info) SnvmapAndMask(outMap []bool) {
	if len(outMap) != len(s.SnvMap) {
		log.Panicf("[Error] Trying to calculate (SnvMap AND Mask), but ")
		log.Panicf("the length of output []boole vector is not equal to reference length\n")
	}
	for i, b := range s.SnvMap {
		if s.ShowMask[i] && b {
			outMap[i] = true
		} else {
			outMap[i] = false
		}
	}
}

// IsAlphabet tell you where c is a alphbet
func IsAlphabet(c byte) bool {
	if 65 <= c && c <= 122 {
		if c <= 90 || c >= 97 {
			return true
		}
	}
	return false
}

// IsEqualAlphabet ignore and tell you where two character is same alphabet
func IsEqualAlphabet(c, d byte) bool {
	if (c-d)%32 == 0 {
		if IsAlphabet(c) && IsAlphabet(d) {
			return true
		}
	}
	return false
}
