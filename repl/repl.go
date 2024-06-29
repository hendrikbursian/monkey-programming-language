package repl

import (
	"bufio"
	"fmt"
	"github.com/hendrikbursian/monkey-programming-language/evaluator"
	"github.com/hendrikbursian/monkey-programming-language/lexer"
	"github.com/hendrikbursian/monkey-programming-language/object"
	"github.com/hendrikbursian/monkey-programming-language/parser"
	"io"
	"os"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			os.Exit(0)
		}

		lexer := lexer.New(line)
		parser := parser.New(lexer)

		program := parser.ParseProgram()

		if len(parser.Errors()) != 0 {
			printParserErrors(out, parser.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

const MONKEY_FACE = `
            __,__
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

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Whoops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errros:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
