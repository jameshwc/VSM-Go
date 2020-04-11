package main

import (
	"log"
	"math"
)

type key struct {
	i, j int
}
type Sparse struct {
	r    int
	c    int
	size int
	// rows []int
	// cols []int
	// data []float64
	elements map[key]float64
}

func NewSparse(r int, c int, size int) *Sparse {
	if uint(r) < 0 {
		log.Fatal("Row Index is Invalid")
	}
	if uint(c) < 0 {
		log.Fatal("Column Index is Invalid")
	}
	// s := &Sparse{r: r, c: c, rows: make([]int, size), cols: make([]int, size), data: make([]float64, size), size: 0}
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
	// s.rows[s.size] = i
	// s.cols[s.size] = j
	// s.data[s.size] = v
	s.elements[key{i, j}] = v
	s.size++
}

func (s *Sparse) Get(r, c int) float64 {
	// for i := 0; i < s.size; i++ {
	// 	if s.rows[i] == r && s.cols[i] == c {
	// 		return s.data[i]
	// 	}
	// }
	// return 0.0
	return s.elements[key{r, c}]
}

func (s *Sparse) PrintRow(r int) {
	// fmt.Println(s.rows, s.cols, s.data)
	// for i := 0; i < s.size; i++ {
	// 	if s.rows[i] == r {
	// 		fmt.Printf("(%d %d) %f\n", r, s.cols[i], s.data[i])
	// 	}
	// }
}

func (s *Sparse) L2Norm() {
	sum := make([]float64, s.r)
	// for i := 0; i < s.size; i++ {
	// 	sum[s.rows[i]] += s.data[i] * s.data[i]
	// }
	for k, v := range s.elements {
		sum[k.i] += v * v
	}
	for i := range sum {
		sum[i] = math.Sqrt(sum[i])
	}
	// for i := 0; i < s.size; i++ {
	// 	s.data[i] = s.data[i] / sum[s.rows[i]]
	// }
	for k, v := range s.elements {
		s.elements[k] = v / sum[k.i]
	}
}
