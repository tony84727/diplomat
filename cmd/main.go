package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MinecraftXwinP/diplomat"
)

func printUsage() {
	fmt.Println("diplomat [translation file]")
}

func getTranslationFile() string {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(-1)
	}
	return os.Args[1]
}

func main() {
	d, err := diplomat.NewDiplomatForFile(getTranslationFile(), "out")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", d.GetOutline())
}
