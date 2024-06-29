package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/hendrikbursian/monkey-programming-language/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	name := strings.ToUpper(user.Username[:1]) + user.Username[1:]

	fmt.Printf("Hello %s! This is the Monkey programming language! \n", name)
	fmt.Printf("Type commands here!\n\n")
	repl.Start(os.Stdin, os.Stdout)
}
