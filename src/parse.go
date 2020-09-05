package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func parseVocab(filename string) map[rune]int {
	vocabID := make(map[rune]int)
	vocabFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("read %s", filename)
	}
	defer vocabFile.Close()
	vocabScanner := bufio.NewScanner(vocabFile)
	vocabScanner.Split(bufio.ScanLines)
	for i := 0; vocabScanner.Scan(); i++ {
		vocabID[[]rune(vocabScanner.Text())[0]] = i
	}
	return vocabID
}

func parseFileList(filename string) (ID2fileName []string) {
	fileListFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("read %s", filename)
	}
	defer fileListFile.Close()
	fileScanner := bufio.NewScanner(fileListFile)
	fileScanner.Split(bufio.ScanLines)
	for i := 0; fileScanner.Scan(); i++ {
		fileName := filepath.Base(strings.ToLower(fileScanner.Text()))
		ID2fileName = append(ID2fileName, fileName)
	}
	return
}

func parseInverted(filename string, fileNum int) (gramID map[gram]int, IDF []float64, docsLen []int, docSum int, matrixSize int) {
	gramID = make(map[gram]int)
	docsLen = make([]int, fileNum)
	invertedFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("read %s", filename)
	}
	defer invertedFile.Close()
	invertedScanner := bufio.NewScanner(invertedFile)
	invertedScanner.Split(bufio.ScanLines)
	for i := 0; invertedScanner.Scan(); i++ {
		splits := strings.SplitN(invertedScanner.Text(), " ", 3)
		vocab1, _ := strconv.Atoi(splits[0])
		vocab2, _ := strconv.Atoi(splits[1])
		n, _ := strconv.Atoi(splits[2])
		g := gram{vocab1, vocab2}
		gramID[g] = i
		matrixSize += n
		for j := 0; j < n; j++ {
			invertedScanner.Scan()
			splits := strings.SplitN(invertedScanner.Text(), " ", 2)
			docID, _ := strconv.Atoi(splits[0])
			freq, _ := strconv.Atoi(splits[1])
			docsLen[docID] += freq
			docSum += freq
		}
		IDF = append(IDF, math.Log(float64(fileNum+1)/float64(n+1))+1)
	}
	return
}

func genTermFrequency(invertFileName string, docsLen []int, avgLen float64, IDF []float64, fileNum int, gramNum int, matrixSize int) *Sparse {
	invertedFile, err := os.Open(invertFileName)
	if err != nil {
		log.Fatalf("read %s", invertFileName)
	}
	defer invertedFile.Close()
	termFrequency := NewSparse(fileNum, gramNum, matrixSize)
	invertedScanner := bufio.NewScanner(invertedFile)
	invertedScanner.Split(bufio.ScanLines)
	bar := progressbar.NewOptions(gramNum, progressbar.OptionSetWriter(os.Stderr))

	for i := 0; invertedScanner.Scan(); i++ {
		n, _ := strconv.Atoi(strings.SplitN(invertedScanner.Text(), " ", 3)[2])
		for j := 0; j < n; j++ {
			invertedScanner.Scan()
			splits := strings.SplitN(invertedScanner.Text(), " ", 2)
			docID, _ := strconv.Atoi(splits[0])
			freq, _ := strconv.ParseFloat(splits[1], 64)
			val := (okapi + 1) * freq / (okapi + freq)
			normalize := (1 - normB) + normB*float64(docsLen[docID])/avgLen
			val = val / normalize * IDF[i]
			termFrequency.Set(docID, i, val)
		}
		bar.Add(1)
	}
	termFrequency.L2Norm()
	return termFrequency
}

func parse(modelDir string) data {
	fmt.Fprintln(os.Stderr, "Parsing...")
	vocabID := parseVocab(filepath.Join(modelDir, "vocab.all"))
	ID2fileName := parseFileList(filepath.Join(modelDir, "file-list"))
	fileNum := len(ID2fileName)
	gramID, IDF, docsLen, docSum, matrixSize := parseInverted(filepath.Join(modelDir, "inverted-file"), fileNum)
	gramNum := len(gramID)
	avgLen := float64(docSum) / float64(fileNum)
	fmt.Fprintln(os.Stderr, "Parsing finished... Now creating tf-idf...")
	termFrequency := genTermFrequency(filepath.Join(modelDir, "inverted-file"), docsLen, avgLen, IDF, fileNum, gramNum, matrixSize)
	return data{vocabID: vocabID, ID2fileName: ID2fileName, gramID: gramID, termFrequency: termFrequency}
}
