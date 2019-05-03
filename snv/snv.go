package snv

import (
	"log"
)

// UpdateSNVMapShowMapSiteCount combine some logical to boost
func UpdateSNVMapShowMapSiteCount(ref []byte, seq []byte, maskChar []byte, snvMap []bool, showMap []bool, counter map[byte][]uint8) {
	if len(ref) < len(seq) {
		log.Panicf("[Error] reference length is shoter!\n")
	}
	if len(snvMap) != len(ref) || len(showMap) != len(ref) {
		log.Panicf("[Error] map size is not equal to ref length!\n")
	}
	for _, ba := range counter {
		if len(ba) != len(ref) {
			log.Panicf("[Error] counter map size is not equal to ref length!\n")
		}
	}
	for i, nt := range seq {
		if !IsEqualAlphabet(nt, ref[i]) {
			snvMap[i] = true
		}
		for _, ex := range maskChar {
			if IsEqualAlphabet(nt, ex) {
				showMap[i] = false
			}
		}
		for c, arr := range counter {
			if IsEqualAlphabet(c, nt) {
				arr[i]++
			}
		}
	}
}

// UpdateSNVMap update snv
func UpdateSNVMap(snvMap []bool, ref []byte, seq []byte) {
	if len(ref) < len(seq) {
		log.Panicf("[Error] reference length is shoter!\n")
	}
	for i, nt := range seq {
		if !IsEqualAlphabet(nt, ref[i]) {
			snvMap[i] = true
		}
	}
}

// UpdateShowMapByMaskChar update...
func UpdateShowMapByMaskChar(showMap []bool, seq []byte, maskChar []byte) {
	// Panic if seq length > reference length
	if len(seq) > len(showMap) {
		log.Panicf("[Error] Length of sequence bigger then showMap.\n")
	}
	for i, c := range seq {
		// Accumulate exlude Site
		for _, ex := range maskChar {
			if IsEqualAlphabet(c, ex) {
				showMap[i] = false
			}
		}
	}
}

// UpdateSiteNtCount update site nt count
func UpdateSiteNtCount(counter map[byte][]uint8, seq []byte) {
	for _, a := range counter {
		if len(a) < len(seq) {
			log.Panicf("[Error] Sequence is longer than nt counter: %d > %d\n", len(seq), len(a))
		}
	}
	for i, nt := range seq {
		for c, arr := range counter {
			if c == nt {
				arr[i]++
			}
		}
	}
}

// BoolArrAND and operate
func BoolArrAND(out, first, second []bool) {
	if len(first) != len(second) || len(first) != len(out) {
		log.Panicf("[Error] Three array length is not equal.\n")
	}
	for i, b := range first {
		if b && second[i] {
			out[i] = true
		} else {
			out[i] = false
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
