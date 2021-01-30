package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const PROMPT = ">> "

func StartConsole(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)

		program := p.ParseProgram()

		if len(p.GetErrors()) != 0 {
			PrintParserErrors(out, p.GetErrors())
			continue
		}

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect() + "\n")
		}
	}
}

const MonkeyFace = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func PrintParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MonkeyFace)
	io.WriteString(out, "Woops! We run into some monkey business here!\n")
	io.WriteString(out, " parser error:\n")

	for _, msg := range errors {
		io.WriteString(out, "\t" + msg + "\n")
	}
}