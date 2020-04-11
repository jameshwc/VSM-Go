package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

type query struct {
	Number    string `xml:"number"`
	Title     string `xml:"title"`
	Question  string `xml:"question"`
	Narrative string `xml:"narrative"`
	Concepts  string `xml:"concepts"`
	Weight    []float64
}

type queries struct {
	Q               []query `xml:"topic"`
	num             int
	titleWeight     float64
	questionWeight  float64
	conceptWeight   float64
	narrativeWeight float64
}

func parseQuery(queryFilePath string, gramNum int, titleWeight, questionWeight, conceptWeight, narrativeWeight float64) queries {
	fmt.Fprintln(os.Stderr, "Parsing queries...")
	queryFile, err := ioutil.ReadFile(queryFilePath)
	if err != nil {
		log.Fatal("read query file")
	}
	q := queries{}
	err = xml.Unmarshal(queryFile, &q)
	if err != nil {
		log.Fatal("parse xml: ", err)
	}
	for i := range q.Q {
		if strings.HasPrefix(q.Q[i].Question, "\n查詢") {
			q.Q[i].Question = q.Q[i].Question[7:] // rune take 3 char and newline take one
		}
		if strings.HasPrefix(q.Q[i].Narrative, "\n相關文件內容") {
			q.Q[i].Narrative = q.Q[i].Narrative[19:]
		}
		if strings.HasPrefix(q.Q[i].Narrative, "應") {
			q.Q[i].Narrative = q.Q[i].Narrative[3:]
		}
		if strings.HasPrefix(q.Q[i].Narrative, "包括") {
			q.Q[i].Narrative = q.Q[i].Narrative[6:]
		}
		if strings.HasPrefix(q.Q[i].Narrative, "說明") {
			q.Q[i].Narrative = q.Q[i].Narrative[6:]
		}
		if strings.HasSuffix(q.Q[i].Narrative, "不相關的。\n") || strings.HasSuffix(q.Q[i].Narrative, "不相關。\n") {
			sp := strings.Split(q.Q[i].Narrative, "。")
			q.Q[i].Narrative = strings.Join(sp[:len(sp)-2], "")
		}
		q.Q[i].Weight = make([]float64, gramNum)
	}
	q.titleWeight = titleWeight
	q.conceptWeight = conceptWeight
	q.narrativeWeight = narrativeWeight
	q.questionWeight = questionWeight
	q.num = len(q.Q)
	fmt.Fprintln(os.Stderr, "Parsing queries finished... Now calculate the weight of queries...")
	return q
}

func (q *queries) calcWeight(gramID map[gram]int, vocabID map[rune]int) {
	bar := progressbar.NewOptions(q.num, progressbar.OptionSetWriter(os.Stderr))
	for i := range q.Q {
		var prevChar rune
		iterString := func(target string, w float64) {
			for idx, char := range target {
				if v, ok := gramID[gram{vocabID[char], -1}]; ok {
					q.Q[i].Weight[v] += w
				}
				if idx != 0 {
					if v, ok := gramID[gram{vocabID[prevChar], vocabID[char]}]; ok {
						q.Q[i].Weight[v] += w
					}
				}
				prevChar = char
			}
		}
		iterString(q.Q[i].Title, q.titleWeight)
		iterString(q.Q[i].Concepts, q.conceptWeight)
		iterString(q.Q[i].Narrative, q.narrativeWeight)
		iterString(q.Q[i].Question, q.questionWeight)
		bar.Add(1)
	}
}
