package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("Недостаточно аргументов")
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
