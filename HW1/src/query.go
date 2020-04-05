package main

import (
	"log"
	"os"
)

type query struct {
	number, title, question, narrative, concepts string
}

func parseQuery(queryFilePath string) {
	queryFile, err := os.Open(queryFilePath)
	if err != nil {
		log.Fatal("read query file")
	}
	defer queryFile.Close()

}
