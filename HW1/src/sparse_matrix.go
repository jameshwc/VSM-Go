package main

import (
	"log"
	"math"
)

type key struct {
	i, j int
}
type Sparse struct {
	r        int
	c        int
	size     int
	elements map[key]float64
}

func NewSparse(r int, c int, size int) *Sparse {
	if uint(r) < 0 {
		log.Fatal("Row Index is Invalid")
	}
	if uint(c) < 0 {
		log.Fatal("Column Index is Invalid")
	}
	s := &Sparse{r: r, c: c, size: 0, elements: make(map[key]float64, size)}
	return s
}

func (s *Sparse) Set(i, j int, v float64) {
	if uint(i) < 0 || uint(i) >= uint(s.r) {
		log.Fatal("Row Index Out of Range")
	}
	if uint(j) < 0 || uint(j) >= uint(s.c) {
		log.Fatal("Column Index Out of Range")
	}
	s.elements[key{i, j}] = v
	s.size++
}

func (s *Sparse) Get(r, c int) float64 {
	return s.elements[key{r, c}]
}

func (s *Sparse) L2Norm() {
	sum := make([]float64, s.r)
	for k, v := range s.elements {
		sum[k.i] += v * v
	}
	for i := range sum {
		sum[i] = math.Sqrt(sum[i])
	}
	for k, v := range s.elements {
		s.elements[k] = v / sum[k.i]
	}
}
