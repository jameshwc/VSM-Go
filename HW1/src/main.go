package main

import (
	"flag"
)

type gram struct {
	vocab1, vocab2 int
}

type data struct {
	vocabID       map[rune]int
	fileID        map[string]int
	gramID        map[gram]int
	ID2fileName   []string
	termFrequency *Sparse
	IDF           []float64
	docsLen       []int
	docSum        int
}

const (
	okapi           = 1.5
	normB           = 0.5
	titleWeight     = 5
	questionWeight  = 4
	conceptWeight   = 15
	narrativeWeight = 2
	maxRetrieveNum  = 60
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
	dat := parse(modelDir, okapi, normB)
	q := parseQuery(queryFilePath, len(dat.gramID), titleWeight, questionWeight, conceptWeight, narrativeWeight)
	q.calcWeight(dat.gramID, dat.vocabID)
	result := newPredicts(q.num, maxRetrieveNum, len(dat.fileID))
	result.predict(dat.termFrequency, q.Q, dat.ID2fileName)
	result.output(rankListPath, q.num)
}
