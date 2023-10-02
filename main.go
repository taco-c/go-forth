package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	isDebugging = true
	debugWidth  = 26
)

func main() {
	var scanner = bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	var isComment = false
	var isFuncScope = false

	for scanner.Scan() {
		var word = scanner.Text()
		if word == ")" {
			debug(word)
			fmt.Print("*/")
			isComment = false
		} else if isComment {
			debug(word)
			fmt.Printf("%s", word)
		} else if word == ":" {
			scanner.Scan()
			debug(fmt.Sprintf("%s %s", word, scanner.Text()))
			fmt.Printf("func %s() {\n", scanner.Text())
			debug("")
			fmt.Printf("\tvar funcStack = make([]interface{}, 0)")
			isFuncScope = true
		} else if word == ";" {
			debug(word)
			fmt.Print("}")
			isFuncScope = false
		} else if word == "(" {
			debug(word)
			fmt.Print("/*")
			isComment = true
		} else if word[0] == '"' {
			if !isFuncScope {
				debug(word)
				fmt.Println()
				continue
			}
			var localString = []string{word}
			if word[len(word)-1] != '"' {
				for scanner.Scan() {
					localString = append(localString, scanner.Text())
					if scanner.Text()[len(scanner.Text())-1] == '"' {
						break
					}
				}
			}
			debug(strings.Join(localString, " "))
			fmt.Printf("\tfuncStack = append(funcStack, %s)", strings.Join(localString, " "))
		} else if word == "package" || word == "import" {
			scanner.Scan()
			debug(fmt.Sprintf("%s %s", word, scanner.Text()))
			fmt.Printf(fmt.Sprintf("%s %s", word, scanner.Text()))
		} else {
			debug(word)
			if isFuncScope == true {
				fmt.Printf("\t%s(pop(funcStack))", word)
			}
		}
		fmt.Println()
	}

	fmt.Println(`
func popN(stack []string, n int) []string {
	var val = len(stack) - n
	var value = stack[val:]
	stack = stack[:val]
	return value
}

func pop[T any](s []T) T {
	var value = s[len(s)-1]
	s = s[:len(s)-1]
	return value
}`)
}

func debug(s string) {
	if isDebugging {
		var fmtStr = fmt.Sprintf("$> %%-%ds | ", debugWidth)
		fmt.Fprintf(os.Stderr, fmtStr, s)
	}
}
