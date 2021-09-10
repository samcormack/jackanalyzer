package jacktokenizer

import (
	"bufio"
	"os"
	"strconv"
	s "strings"
	"unicode"
)

type Token int

const (
	KEYWORD Token = iota
	SYMBOL
	IDENTIFIER
	INT_CONST
	STRING_CONST
)

type Keyword int

const (
	CLASS Keyword = iota
	METHOD
	FUNCTION
	CONSTRUCTOR
	INT
	BOOLEAN
	CHAR
	VOID
	VAR
	STATIC
	FIELD
	LET
	DO
	IF
	ELSE
	WHILE
	RETURN
	TRUE
	FALSE
	NULL
	THIS
)

var Keywords = map[string]Keyword{
	"class":       CLASS,
	"constructor": CONSTRUCTOR,
	"function":    FUNCTION,
	"method":      METHOD,
	"field":       FIELD,
	"static":      STATIC,
	"var":         VAR,
	"int":         INT,
	"char":        CHAR,
	"boolean":     BOOLEAN,
	"void":        VOID,
	"true":        TRUE,
	"false":       FALSE,
	"null":        NULL,
	"this":        THIS,
	"let":         LET,
	"do":          DO,
	"if":          IF,
	"else":        ELSE,
	"while":       WHILE,
	"return":      RETURN,
}

var symbols = map[rune]bool{
	'}': true,
	'{': true,
	'(': true,
	')': true,
	'[': true,
	']': true,
	'.': true,
	',': true,
	';': true,
	'+': true,
	'-': true,
	'*': true,
	'/': true,
	'&': true,
	'|': true,
	'<': true,
	'>': true,
	'=': true,
	'~': true,
}

type Tokenizer struct {
	file        *os.File
	reader      *bufio.Reader
	currentRune rune
	tokenType   Token
	tokenText   string
}

func NewTokenizer(file *os.File) *Tokenizer {
	reader := bufio.NewReader(file)
	t := Tokenizer{file: file, reader: reader}
	return &t
}

func (t *Tokenizer) HasMoreTokens() bool {
	var err error
	t.currentRune, _, err = t.reader.ReadRune()

	// Skip whitespace
	for unicode.IsSpace(t.currentRune) {
		t.currentRune, _, err = t.reader.ReadRune()
	}

	// Skip a comment then advance again
	if t.currentRune == '/' {
		t.skipComments()
		return t.HasMoreTokens()
	}

	return err == nil
}

func (t *Tokenizer) skipComments() {
	next, err := t.reader.Peek(1)
	nextRune := rune(next[0])
	if err == nil {
		switch nextRune {
		case '/':
			t.reader.ReadRune()
			t.currentRune, _, _ = t.reader.ReadRune()
			for t.currentRune != '\n' {
				t.currentRune, _, _ = t.reader.ReadRune()

			}

		case '*':
			t.reader.ReadRune()
			r1, _, _ := t.reader.ReadRune()
			t.currentRune, _, _ = t.reader.ReadRune()
			for !(r1 == '*' && t.currentRune == '/') {
				r1 = t.currentRune
				t.currentRune, _, _ = t.reader.ReadRune()

			}

		}
	}
}

func (t *Tokenizer) Advance() {
	var err error
	tokenNotFound := true
	for tokenNotFound {
		tokenNotFound = false
		if isSymbol(t.currentRune) {
			t.tokenType = SYMBOL
		} else if unicode.IsDigit(t.currentRune) {
			t.tokenType = INT_CONST
			t.scanInt()
		} else if isQuote(t.currentRune) {
			t.tokenType = STRING_CONST
			t.scanString()
		} else if isIdentStart(t.currentRune) {
			t.scanWord()
		} else {
			tokenNotFound = true
			t.currentRune, _, err = t.reader.ReadRune()
			if err != nil {
				break
			}
		}
	}
}

func isSymbol(r rune) bool {
	return symbols[r]
}

func (t *Tokenizer) scanInt() {
	var intBuilder s.Builder
	intBuilder.WriteRune(t.currentRune)
	for {
		t.currentRune, _, _ = t.reader.ReadRune()
		if unicode.IsDigit(t.currentRune) {
			intBuilder.WriteRune(t.currentRune)
		} else {
			t.reader.UnreadRune()
			break
		}
	}
	t.tokenText = intBuilder.String()
}

func isQuote(text rune) bool {
	return text == '"'
}

func (t *Tokenizer) scanString() {
	var strBuilder s.Builder
	for {
		t.currentRune, _, _ = t.reader.ReadRune()
		if isQuote(t.currentRune) {
			break
		} else {
			strBuilder.WriteRune(t.currentRune)
		}
	}
	t.tokenText = strBuilder.String()
}

func isIdentStart(r rune) bool {
	return unicode.IsLetter(r)
}

func isIdentAny(r rune) bool {
	return unicode.IsDigit(r) || unicode.IsLetter(r) || r == '_'
}

func (t *Tokenizer) scanWord() {
	var wordBuilder s.Builder
	wordBuilder.WriteRune(t.currentRune)
	for {
		t.currentRune, _, _ = t.reader.ReadRune()
		if isIdentAny(t.currentRune) {
			wordBuilder.WriteRune(t.currentRune)
		} else {
			t.reader.UnreadRune()
			break
		}
	}
	t.tokenText = wordBuilder.String()
	_, found := Keywords[t.tokenText]
	if found {
		t.tokenType = KEYWORD
	} else {
		t.tokenType = IDENTIFIER
	}
}

func (t *Tokenizer) TokenType() Token {
	return t.tokenType
}

func (t *Tokenizer) TokenName() string {
	switch t.TokenType() {
	case KEYWORD:
		return "keyword"
	case SYMBOL:
		return "symbol"
	case IDENTIFIER:
		return "identifier"
	case INT_CONST:
		return "intConst"
	case STRING_CONST:
		return "strConst"
	}
	return "err"
}

func (t *Tokenizer) KeyWord() Keyword {
	return Keywords[t.tokenText]
}

func (t *Tokenizer) KeyWordName() string {
	return t.tokenText
}

func (t *Tokenizer) Symbol() rune {
	return t.currentRune
}

func (t *Tokenizer) Identifier() string {
	return t.tokenText
}

func (t *Tokenizer) IntVal() int {
	intVal, _ := strconv.Atoi(t.tokenText)
	return intVal
}

func (t *Tokenizer) StringVal() string {
	return t.tokenText
}
