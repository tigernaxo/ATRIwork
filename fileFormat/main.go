package fileformat

// Feature 保存特徵範圍、正負股、名稱
type Feature struct {
	Start  int    // gff col 4
	End    int    // gff col 5
	Name   string // gff col9 extract Name=name
	Strand byte   // gff col 7
}

// FeatureSet 保存同一類特徵、列表
type FeatureSet struct {
	Class    string // gff col 3
	Features []*Feature
}
