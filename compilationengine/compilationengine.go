package compilationengine

import (
	"bufio"
	"jackanalyzer/jacktokenizer"
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
	if engine.tokenizer.HasMoreTokens() {
		engine.tokenizer.Advance()

		if engine.tokenizer.TokenType() == jacktokenizer.KEYWORD &&
			engine.tokenizer.KeyWord() == jacktokenizer.CLASS {
			engine.writer.WriteString("Class")
		}
	}
}

func (engine *CompilationEngine) compileClassVarDec() {}

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
