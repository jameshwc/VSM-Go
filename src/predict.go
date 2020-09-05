package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/schollz/progressbar/v3"
)

type sim struct {
	idx int
	dat float64
}

type sims []sim

type predict struct {
	rank []string
	sims sims
}

type predicts []predict

func (s sims) Len() int           { return len(s) }
func (s sims) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sims) Less(i, j int) bool { return s[i].dat > s[j].dat }

func (p *predict) sort(ID2fileName []string) {
	for i := range p.sims {
		p.sims[i].idx = i
	}
	sort.Sort(p.sims)
	n := len(p.rank)
	for i := 0; i < n; i++ {
		p.rank[i] = ID2fileName[p.sims[i].idx]
	}
}

func newPredicts(queryNum, maxRetrieNum, fileNum int) predicts {
	p := make([]predict, queryNum)
	for i := range p {
		p[i].rank = make([]string, maxRetrieNum)
		p[i].sims = make([]sim, fileNum)
	}
	return p
}

func (p predicts) predict(docWeight *Sparse, queries []query, ID2fileName []string) {
	docNum := len(ID2fileName)
	fmt.Fprintln(os.Stderr, "\nNow calc the final result...\n")
	bar := progressbar.NewOptions(len(queries), progressbar.OptionSetWriter(os.Stderr))
	for q := range queries {
		for gid := range queries[q].Weight {
			if queries[q].Weight[gid] != 0.0 {
				for docid := 0; docid < docNum; docid++ {
					p[q].sims[docid].dat += queries[q].Weight[gid] * docWeight.Get(docid, gid)
				}
			}
		}
		p[q].sort(ID2fileName)
		bar.Add(1)
	}
	fmt.Fprintln(os.Stderr, "\nFinish calculating!\n")
}

func (p predicts) rocchio(queries []query, docWeight *Sparse) {
	sumTopNdoc := func(p predict, gid int) float64 {
		s := 0.0
		for i := 0; i < topNumRel; i++ {
			s += docWeight.Get(p.sims[i].idx, gid)
		}
		return s
	}
	sumLastNdoc := func(p predict, gid int) float64 {
		s := 0.0
		for i := 0; i < lastNumRel; i++ {
			s += docWeight.Get(p.sims[maxRetrieveNum-1-i].idx, gid)
		}
		return s
	}
	for q := range queries {
		for gid := range queries[q].Weight {
			if queries[q].Weight[gid] != 0.0 {
				queries[q].Weight[gid] = alpha*queries[q].Weight[gid] + beta/topNumRel*sumTopNdoc(p[q], gid) - c*sumLastNdoc(p[q], gid)
			}
		}
	}
	for predict := range p {
		for s := range p[predict].sims {
			p[predict].sims[s].dat = 0.0
		}
	}
}
func (p predicts) output(CSVpath string, queryNum int) {
	csvfile, err := os.Create(CSVpath)
	if err != nil {
		log.Fatal("Couldn't create the csv file", err)
	}
	defer csvfile.Close()
	csvfile.WriteString("query_id,retrieved_docs\n")
	for i := range p {
		str := fmt.Sprintf("%.3d,%s\n", i+11, strings.Join(p[i].rank, " "))
		csvfile.WriteString(str)
	}
}
