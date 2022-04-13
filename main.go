package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"path/filepath"
	"regexp"
	//"compress/gzip"
	gzip "github.com/klauspost/pgzip"
	"github.com/liserjrqlxue/goUtil/textUtil"

	"github.com/liserjrqlxue/goUtil/osUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
)

// flag
var (
	in     = flag.String("in", "", "input file")
	inList = flag.String("inList", "", "input file list")
	outDir = flag.String("out", "", "output directory")
	offset = flag.Int("offset", 33, "quality offset")
)

// regexp
var (
	isGz = regexp.MustCompile(`\.gz$`)
)

var (
	qualityBins map[byte]byte
)

func main() {
	flag.Parse()
	if *in == "" && *inList == "" || *outDir == "" {
		flag.Usage()
		return
	}

	// init
	// quality to bins
	qualityBins = make(map[byte]byte)
	var bins = [9]int{0, 0, 3, 11, 20, 23, 30, 37, 100}
	for i, bin := range bins {
		bins[i] = bin + *offset
	}
	for i := 0; i < len(bins)-2; i += 2 {
		for j := bins[i]; j < bins[i+2]; j++ {
			qualityBins[byte(j)] = byte(bins[i+1])
		}
	}

	// prepare input
	var (
		inputs []string
	)
	if *in != "" {
		inputs = append(inputs, *in)
	}
	if *inList != "" {
		for _, s := range textUtil.File2Array(*inList) {
			inputs = append(inputs, s)
		}
	}

	simpleUtil.CheckErr(os.MkdirAll(*outDir, 0755))

	for _, input := range inputs {
		var output = filepath.Join(*outDir, filepath.Base(input))
		qualityBin(input, output)
	}
}

func qualityBin(input, output string) {
	var inFile = osUtil.Open(input)
	var outFile = osUtil.Create(output)
	if isGz.MatchString(input) {
		var gr = simpleUtil.HandleError(gzip.NewReader(inFile)).(*gzip.Reader)
		var gw = gzip.NewWriter(outFile)
		quality2bin(gr, gw)
		simpleUtil.CheckErr(inFile.Close())
		simpleUtil.CheckErr(outFile.Close())
	} else {
		quality2bin(inFile, outFile)
	}

}

func quality2bin(input io.ReadCloser, output io.WriteCloser) {
	var n = 0
	var scanner = bufio.NewScanner(input)
	for scanner.Scan() {
		var line = scanner.Bytes()
		n++
		if n%4 == 0 {
			for i, b := range line {
				line[i] = qualityBins[b]
			}
		}
		simpleUtil.HandleError(output.Write(append(line, '\n')))
	}
	simpleUtil.CheckErr(scanner.Err())
	simpleUtil.CheckErr(input.Close())
	simpleUtil.CheckErr(output.Close())
}
