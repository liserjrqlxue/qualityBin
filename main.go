package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"regexp"

	"github.com/liserjrqlxue/goUtil/osUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
)

// flag
var (
	in  = flag.String("in", "", "input file")
	out = flag.String("out", "", "output file")
)

// regexp
var (
	isGz = regexp.MustCompile(`\.gz$`)
)

func main() {
	flag.Parse()
	if *in == "" || *out == "" {
		flag.PrintDefaults()
		return
	}

	// quality to bins
	var qualityBins = make(map[byte]byte)
	var bins = [9]int{0, 0, 3, 11, 20, 23, 30, 37, 100}
	for i := 0; i < len(bins)-1; i += 2 {
		for j := bins[i]; j < bins[i+2]; j++ {
			qualityBins[byte(j)] = byte(bins[i+1])
		}
	}

	// read file
	var input = osUtil.Open(*in)
	defer simpleUtil.DeferClose(input)
	var scanner *bufio.Scanner
	if isGz.MatchString(*in) {
		var gr = simpleUtil.HandleError(gzip.NewReader(input)).(*gzip.Reader)
		scanner = bufio.NewScanner(gr)
	} else {
		scanner = bufio.NewScanner(input)
	}

	// write file
	var output = osUtil.Create(*out)
	defer simpleUtil.DeferClose(output)
	var outZw *gzip.Writer
	if isGz.MatchString(*out) {
		outZw = gzip.NewWriter(output)
		defer simpleUtil.DeferClose(outZw)
	}

	// write
	var n = 0
	for scanner.Scan() {
		var line = scanner.Bytes()
		n++
		if n%4 == 0 {
			for i, b := range line {
				line[i] = qualityBins[b]
			}
		}
		if outZw != nil {
			simpleUtil.HandleError(outZw.Write(append(line, '\n')))
		} else {
			simpleUtil.HandleError(outZw.Write(append(line, '\n')))
		}
	}
}
