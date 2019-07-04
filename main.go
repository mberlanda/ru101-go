package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mberlanda/ru101-go/utils"
)

func main() {
	seedFn := flag.String("seed", "", "a file path where seeds are located")
	flag.Parse()

	if *seedFn != "" {
		log.Println("Started")
		utils.LoadFromFile(*seedFn)
		log.Println("Completed")
		return
	}

	fmt.Println("Usage")
	flag.PrintDefaults()
	os.Exit(1)
}
