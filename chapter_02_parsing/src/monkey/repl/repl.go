package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const PROMPT = ">> "

func StartConsole(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()

		if (!scanned) {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)

		for t := l.NextToken(); t.Type != token.EOF; t = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", t)
		}
	}
}