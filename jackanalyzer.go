package main

import (
	"bufio"
	"jackanalyzer/jacktokenizer"
	"log"
	"os"
	"path/filepath"
	"strconv"
	s "strings"
)

func main() {
	// Get input file and open
	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	fname := s.TrimSuffix(arg, ".jack")
	var outfileName string
	if fname == arg {
		outfileName = filepath.Join(fname, filepath.Base(fname)) + "T.xml"
	} else {
		outfileName = fname + "T.xml"
	}

	// Create output file and writer
	outfile, err := os.Create(outfileName)
	check(err)
	defer outfile.Close()

	infile, err := os.Open(arg)
	check(err)
	defer infile.Close()

	testTokenizer(infile, outfile)

}

func testTokenizer(infile *os.File, outfile *os.File) {
	writer := bufio.NewWriter(outfile)
	defer writer.Flush()
	jt := jacktokenizer.NewTokenizer(infile)
	for jt.HasMoreTokens() {
		jt.Advance()
		writer.WriteString("<" + jt.TokenName() + ">")
		switch jt.TokenType() {
		case jacktokenizer.KEYWORD:
			writer.WriteString(jt.KeyWordName())
		case jacktokenizer.IDENTIFIER:
			writer.WriteString(jt.Identifier())
		case jacktokenizer.STRING_CONST:
			writer.WriteString(jt.StringVal())
		case jacktokenizer.INT_CONST:
			writer.WriteString(strconv.Itoa(jt.IntVal()))
		case jacktokenizer.SYMBOL:
			writer.WriteRune(jt.Symbol())
		}
		writer.WriteString("</" + jt.TokenName() + ">\n")
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
