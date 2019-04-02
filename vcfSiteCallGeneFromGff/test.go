package main

import (
	"fmt"
	"os"
)

func appendSlice() {
	s := make([]int, 10)
	fmt.Println(len(s))
	fmt.Println(cap(s))
}
func writeFile() {
	f, _ := os.Create("test.tsv")
	defer f.Close()
	f.WriteString("aaa")
}
