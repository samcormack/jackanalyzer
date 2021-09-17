package main

import (
	"bufio"
	"io/ioutil"
	"jackanalyzer/compilationengine"
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

	if fname == arg {
		// Run on files in directory
		files, err := ioutil.ReadDir(fname)
		check(err)
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".jack" {
				analyzeToXML(filepath.Join(arg, file.Name()))
			}
		}
	} else if filepath.Ext(arg) == ".jack" {
		analyzeToXML(filepath.Clean(arg))
	}

}

func analyzeToXML(filename string) {
	outfileName := s.TrimSuffix(filename, filepath.Ext(filename))
	outfileName = outfileName + ".xml"

	outfile, err := os.Create(outfileName)
	check(err)
	defer outfile.Close()

	infile, err := os.Open(filename)
	check(err)
	defer infile.Close()

	engine := compilationengine.NewCompilationEngine(infile, outfile)
	defer engine.Finish()
	engine.CompileClass()
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
