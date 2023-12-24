package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! Have fun with the Monkey programming language!\n", user.Username)
	fmt.Print("Feel free to hack away\n")
	repl.Start(os.Stdin, os.Stdout)
}
