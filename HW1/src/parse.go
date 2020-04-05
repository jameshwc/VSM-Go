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

	"github.com/james-bowman/sparse"
)

type gram struct {
	vocab1, vocab2 int
}

func parse(modelDir string, okapi float64) (vocabID, fileID map[string]int, gramID map[gram]int, termFrequency *sparse.DOK, IDF []float64) {
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
	vocabID = make(map[string]int)
	fileID = make(map[string]int)
	gramID = make(map[gram]int)
	vocabScanner := bufio.NewScanner(vocabFile)
	vocabScanner.Split(bufio.ScanLines)
	for i := 0; vocabScanner.Scan(); i++ {
		vocabID[vocabScanner.Text()] = i
	}
	fileScanner := bufio.NewScanner(fileListFile)
	fileScanner.Split(bufio.ScanLines)
	for i := 0; fileScanner.Scan(); i++ {
		fileID[fileScanner.Text()] = i
	}
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
		for j := 0; j < n; j++ {
			invertedScanner.Scan()
		}
	}

	/* Generate the term-frequency matrix */

	fileNum := len(fileID)
	termFrequency = sparse.NewDOK(fileNum, len(gramID))
	invertedFile.Seek(0, io.SeekStart)
	invertedScanner = bufio.NewScanner(invertedFile)
	invertedScanner.Split(bufio.ScanLines)
	docsLen := make([]int, fileNum)
	docSum := 0
	for i := 0; invertedScanner.Scan(); i++ {
		n, _ := strconv.Atoi(strings.SplitN(invertedScanner.Text(), " ", 3)[2])
		for j := 0; j < n; j++ {
			splits := strings.SplitN(invertedScanner.Text(), " ", 2)
			docID, _ := strconv.Atoi(splits[0])
			freq, _ := strconv.ParseFloat(splits[1], 64)
			// fmt.Println(docID, freq)
			docsLen[docID] += int(freq)
			docSum += int(freq)
			val := (okapi + 1) * freq / (okapi + freq)
			termFrequency.Set(docID, i, val)
			invertedScanner.Scan()
		}
		IDF = append(IDF, math.Log(float64(fileNum+1)/float64(n+1))+1)
	}
	fmt.Print(docsLen, len(docsLen), docSum)
	return
}
