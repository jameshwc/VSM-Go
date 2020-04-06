package main

import (
	"flag"
	"fmt"
)

func main() {
	var queryFilePath, rankListPath, modelDir, NTCIRdir string
	var isRelFeedback bool
	flag.BoolVar(&isRelFeedback, "r", false, "If specified, turn on the relevance feedback in your program.")
	flag.StringVar(&queryFilePath, "i", "queries/query-test.xml", "The input query file.")
	flag.StringVar(&rankListPath, "o", "test.csv", "The output ranked list file.")
	flag.StringVar(&modelDir, "m", "model/", "The input model directory, which includes three files:\n\tmodel-dir/vocab.all\n\tmodel-dir/file-list\n\tmodel-dir/inverted-index")
	flag.StringVar(&NTCIRdir, "d", "CIRB010/", "The directory of NTCIR documents, which is the path name of CIRB010 directory.")
	flag.Parse()
	dat := parse(modelDir, 1.5)
	fmt.Println(dat.docSum, len(dat.fileID))

}
