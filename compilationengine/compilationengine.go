package compilationengine

import (
	"bufio"
	"fmt"
	"jackanalyzer/jacktokenizer"
	"log"
	"os"
)

type CompilationEngine struct {
	tokenizer *jacktokenizer.Tokenizer
	outfile   *os.File
	writer    *bufio.Writer
}

func NewCompilationEngine(infile *os.File, outfile *os.File) *CompilationEngine {
	jt := jacktokenizer.NewTokenizer(infile)
	writer := bufio.NewWriter(outfile)
	return &CompilationEngine{tokenizer: jt, outfile: outfile, writer: writer}
}

func (engine *CompilationEngine) Finish() {
	engine.writer.Flush()
}

func (engine *CompilationEngine) CompileClass() {
	engine.writer.WriteString("<class>\n")
	// Check class keyword
	if engine.tokenizer.HasMoreTokens() {
		engine.tokenizer.Advance()

		if engine.tokenizer.TokenType() == jacktokenizer.KEYWORD &&
			engine.tokenizer.KeyWord() == jacktokenizer.CLASS {
			engine.writer.WriteString("<keyword> class </keyword>\n")
		} else {
			log.Fatal("Did not find expected class keyword.")
		}
	}

	// Get class name
	if engine.tokenizer.HasMoreTokens() {
		engine.tokenizer.Advance()

		if engine.tokenizer.TokenType() == jacktokenizer.IDENTIFIER {
			engine.printIdentifier(engine.tokenizer.Identifier())
		} else {
			log.Fatal("Did not find an identifier for class.")
		}
	}

	// Check {
	if engine.tokenizer.HasMoreTokens() {
		engine.tokenizer.Advance()

		if engine.tokenizer.TokenType() == jacktokenizer.SYMBOL &&
			engine.tokenizer.Symbol() == '{' {
			engine.printSymbol('{')
		} else {
			log.Fatal("Did not find expected { after class declaration.")
		}
	}

	// compile ClassVarDec's

	for engine.tokenizer.HasMoreTokens() {
		engine.tokenizer.Advance()

		if engine.tokenizer.TokenType() == jacktokenizer.KEYWORD &&
			(engine.tokenizer.KeyWord() == jacktokenizer.STATIC || engine.tokenizer.KeyWord() == jacktokenizer.FIELD) {
			engine.compileClassVarDec()
		} else {
			break
		}
	}

	//compile subroutineDec's
	for engine.tokenizer.HasMoreTokens() {
		engine.tokenizer.Advance()

		if engine.tokenizer.TokenType() == jacktokenizer.KEYWORD &&
			(engine.tokenizer.KeyWord() == jacktokenizer.CONSTRUCTOR ||
				engine.tokenizer.KeyWord() == jacktokenizer.FUNCTION ||
				engine.tokenizer.KeyWord() == jacktokenizer.METHOD) {
			engine.compileSubroutine()
		} else {
			break
		}
	}

	// Check }
	if engine.tokenizer.HasMoreTokens() {
		engine.tokenizer.Advance()

		if engine.tokenizer.TokenType() == jacktokenizer.SYMBOL &&
			engine.tokenizer.Symbol() == '}' {
			engine.printSymbol('}')
		} else {
			engine.writer.Flush()
			log.Fatal("Did not find expected } after class declaration.")
		}
	}

	engine.writer.WriteString("</class>")

}

func (engine *CompilationEngine) compileClassVarDec() {
	engine.writer.WriteString("<classVarDec>")
	engine.printKeyword(engine.tokenizer.KeyWordName())

	// Type dec
	if engine.tokenizer.HasMoreTokens() {
		engine.tokenizer.Advance()
		if engine.isType() {
			engine.printKeyword(engine.tokenizer.KeyWordName())
		} else {
			log.Fatal("Did not find type declaration for class variable.")
		}
	}

	engine.writer.WriteString("</classVarDec>")
}

func (engine *CompilationEngine) compileSubroutine() {}

func (engine *CompilationEngine) compileParameterList() {}

func (engine *CompilationEngine) compileVarDec() {}

func (engine *CompilationEngine) compileStatements() {}

func (engine *CompilationEngine) compileDo() {}

func (engine *CompilationEngine) compileLet() {}

func (engine *CompilationEngine) compileWhile() {}

func (engine *CompilationEngine) compileReturn() {}

func (engine *CompilationEngine) compileIf() {}

func (engine *CompilationEngine) compileExpression() {}

func (engine *CompilationEngine) compileTerm() {}

func (engine *CompilationEngine) compileExpressionList() {}

func (engine *CompilationEngine) isType() bool {
	if engine.tokenizer.TokenType() == jacktokenizer.KEYWORD {
		return engine.tokenizer.KeyWord() == jacktokenizer.INT ||
			engine.tokenizer.KeyWord() == jacktokenizer.CHAR ||
			engine.tokenizer.KeyWord() == jacktokenizer.BOOLEAN
	} else {
		return engine.tokenizer.TokenType() == jacktokenizer.IDENTIFIER
	}

}

func (engine *CompilationEngine) printSymbol(r rune) {
	engine.writer.WriteString(
		fmt.Sprintf(
			"<symbol> %s </symbol>\n",
			string(r),
		),
	)
}

func (engine *CompilationEngine) printIdentifier(s string) {
	engine.writer.WriteString(
		fmt.Sprintf(
			"<identifier> %s </identifier>\n",
			s,
		),
	)
}

func (engine *CompilationEngine) printKeyword(s string) {
	engine.writer.WriteString(
		fmt.Sprintf(
			"<keyword> %s </keyword>\n",
			s,
		),
	)
}
