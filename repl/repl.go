package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/eval"
	"monkey/lexer"
	"monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		lex := lexer.New(line)
		parser := parser.New(lex)

		errors, program := parser.ParseProgram()

		if errors != nil {
			outputErrors(out, errors)
			continue
		}

		value := eval.Eval(program)

		if value != nil {
			fmt.Fprintf(out, "%+v\n", value.Inspect())
		}
	}
}

func outputErrors(out io.Writer, errors []string) {
	fmt.Fprintf(out, "😅 Ooops ... we encountered some errors:\n")

	for _, error := range errors {
		fmt.Fprintf(out, "\t %s\n", error)
	}
}
