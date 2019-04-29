package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type config struct {
	intensity float64
	gapRatio  float64
	outDir    string
	refSeq    string
	fileList  []string
}

func newConfig(file string) *config {

	var intensity, gapRatio float64
	prefixIntensity := "intensity"
	prefixGapRatio := "gapRatio"
	var refSeq, outDir string
	prefixRefSeq := "reference"
	prefixOutDir := "outdir"

	configStartLine := "[config]"
	fileListStartLine := "[file list]"
	asignSymbol := "="

	fileList := make([]string, 0, 10)
	inFileList, inConfig := false, false

	f, err := os.Open(file)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		switch line := strings.TrimSpace(scanner.Text()); line {
		case fileListStartLine:
			inFileList, inConfig = true, false
		case configStartLine:
			inFileList, inConfig = false, true
		default:
			if inFileList && !inConfig {
				fileList = append(fileList, line)
			} else if !inFileList && inConfig {
				switch strings.Split(line, asignSymbol)[0] {
				case prefixGapRatio:
					// fmt.Printf("[log] parsing %v \n")
					gapRatio, err = strconv.ParseFloat(strings.Trim(line, prefixGapRatio+asignSymbol), 64)
					if err != nil {
						log.Panicf("[Error] Error in Parsing %s to float64 \n%v\n", line, err)
					}
				case prefixIntensity:
					intensity, err = strconv.ParseFloat(strings.Trim(line, prefixIntensity+asignSymbol), 64)
					if err != nil {
						log.Panicf("[Error] Error in Parsing %s to float64 \n%v\n", line, err)
					}
				case prefixOutDir:
					outDir = strings.Trim(line, prefixOutDir+asignSymbol)
				case prefixRefSeq:
					refSeq = strings.Trim(line, prefixRefSeq+asignSymbol)
				default:
				}
			}
		}
	}

	return &config{
		intensity: intensity,
		gapRatio:  gapRatio,
		refSeq:    refSeq,
		outDir:    outDir,
		fileList:  fileList,
	}
}
