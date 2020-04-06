package main

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

type Sparse struct {
	r    int
	c    int
	size int
	rows []int
	cols []int
	data []float64
}

func NewSparse(r int, c int, size int) *Sparse {
	if uint(r) < 0 {
		panic(mat.ErrRowAccess)
	}
	if uint(c) < 0 {
		panic(mat.ErrColAccess)
	}
	s := &Sparse{r: r, c: c, rows: make([]int, size), cols: make([]int, size), data: make([]float64, size), size: 0}

	return s
}
func (s *Sparse) Set(i, j int, v float64) {
	if uint(i) < 0 || uint(i) >= uint(s.r) {
		panic(mat.ErrRowAccess)
	}
	if uint(j) < 0 || uint(j) >= uint(s.c) {
		panic(mat.ErrColAccess)
	}
	s.rows[s.size] = i
	s.cols[s.size] = j
	s.data[s.size] = v
	s.size++
}

func (s *Sparse) PrintRow(r int) {
	// fmt.Println(s.rows, s.cols, s.data)
	for i := 0; i < s.size; i++ {
		if s.rows[i] == r {
			fmt.Printf("(%d %d) %f\n", r, s.cols[i], s.data[i])
		}
	}
}

func (s *Sparse) L2Norm() {
	sum := make([]float64, s.r)
	for i := 0; i < s.size; i++ {
		sum[s.rows[i]] += s.data[i] * s.data[i]
	}
	for i := range sum {
		sum[i] = math.Sqrt(sum[i])
	}
	for i := 0; i < s.size; i++ {
		s.data[i] = s.data[i] / sum[s.rows[i]]
	}
}
