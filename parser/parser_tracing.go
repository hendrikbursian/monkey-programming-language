package parser

import (
	"fmt"
)

const TRACE = false

var traceLevel int = 0

const traceIdentPlaceholder string = "\t"

// func identLevel() string {
// 	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
// }

// func tracePrint(fs string) {
// 	fmt.Printf("%s%s\n", identLevel(), fs)
// }

func tracePrint(fs string) {
	if TRACE {
		fmt.Printf("%*s%s\n", (traceLevel-1)*4, "", fs)
	}
}

func incIdent() { traceLevel = traceLevel + 1 }
func decIdent() { traceLevel = traceLevel - 1 }

func trace(msg string) string {
	incIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	tracePrint("END " + msg)
	decIdent()
}
