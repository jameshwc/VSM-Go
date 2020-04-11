package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func parse(modelDir string, okapi float64, normB float64) data {
	fmt.Fprintln(os.Stderr, "Parsing...")
	vocabFile, err := os.Open(filepath.Join(modelDir, "vocab.all"))
	if err != nil {
		log.Fatal("read vocab.all")
	}
	defer vocabFile.Close()
	fileListFile, err := os.Open(filepath.Join(modelDir, "file-list"))
	if err != nil {
		log.Fatal("read file-list")
	}
	defer fileListFile.Close()
	invertedFile, err := os.Open(filepath.Join(modelDir, "inverted-file"))
	if err != nil {
		log.Fatal("read inverted-file")
	}
	defer invertedFile.Close()
	var ID2fileName []string
	var IDF []float64
	vocabID := make(map[rune]int)
	fileID := make(map[string]int)
	gramID := make(map[gram]int)
	vocabScanner := bufio.NewScanner(vocabFile)
	vocabScanner.Split(bufio.ScanLines)
	for i := 0; vocabScanner.Scan(); i++ {
		vocabID[[]rune(vocabScanner.Text())[0]] = i
	}
	fileScanner := bufio.NewScanner(fileListFile)
	fileScanner.Split(bufio.ScanLines)
	for i := 0; fileScanner.Scan(); i++ {
		fileName := fileScanner.Text()
		ID2fileName = append(ID2fileName, fileName)
		fileID[fileName] = i
	}
	fileNum := len(fileID)
	docsLen := make([]int, fileNum)
	docSum := 0
	matrixSize := 0
	invertedScanner := bufio.NewScanner(invertedFile)
	invertedScanner.Split(bufio.ScanLines)
	for i := 0; invertedScanner.Scan(); i++ {
		splits := strings.SplitN(invertedScanner.Text(), " ", 3)
		vocab1, _ := strconv.Atoi(splits[0])
		vocab2, _ := strconv.Atoi(splits[1])
		g := gram{vocab1, vocab2}
		gramID[g] = i
		n, err := strconv.Atoi(splits[2])
		if err != nil {
			log.Fatal("read inverted-file: n is not a number")
		}
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
	gramNum := len(gramID)
	avgLen := float64(docSum) / float64(fileNum)
	fmt.Fprintln(os.Stderr, "Parsing finished... Now creating tf-idf...")

	/* Generate the term-frequency matrix */

	termFrequency := NewSparse(fileNum, gramNum, matrixSize)
	fmt.Println(termFrequency.r, termFrequency.c)
	invertedFile.Seek(0, io.SeekStart)
	invertedScanner = bufio.NewScanner(invertedFile)
	invertedScanner.Split(bufio.ScanLines)
	bar := progressbar.NewOptions(gramNum, progressbar.OptionSetWriter(os.Stderr))
	for i := 0; invertedScanner.Scan(); i++ {
		n, _ := strconv.Atoi(strings.SplitN(invertedScanner.Text(), " ", 3)[2])
		for j := 0; j < n; j++ {
			invertedScanner.Scan()
			splits := strings.SplitN(invertedScanner.Text(), " ", 2)
			docID, _ := strconv.Atoi(splits[0])
			freq, _ := strconv.ParseFloat(splits[1], 64)
			// fmt.Println(docID, freq)
			val := (okapi + 1) * freq / (okapi + freq)
			normalize := (1 - normB) + normB*float64(docsLen[docID])/avgLen
			val = val / normalize * IDF[i]
			termFrequency.Set(docID, i, val)
		}
		bar.Add(1)
	}
	termFrequency.L2Norm()
	dat := data{vocabID: vocabID, fileID: fileID, ID2fileName: ID2fileName, gramID: gramID, IDF: IDF, docSum: docSum, docsLen: docsLen, termFrequency: termFrequency}
	return dat
}
