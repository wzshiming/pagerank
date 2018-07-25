package pagerank

import (
	"math"
)

// Pagerank is pagerank calculator
type Pagerank struct {
	Matrix   [][]uint64
	Links    []uint64
	Keymap   map[string]uint64
	Rekeymap []string
}

// NewPagerank Create a pagerank calculator
func NewPagerank() *Pagerank {
	return &Pagerank{
		Keymap: map[string]uint64{},
	}
}

// Len returns total number of pages
func (pr *Pagerank) Len() int {
	return len(pr.Rekeymap)
}

// Link mark a link
func (pr *Pagerank) Link(from, to string) {
	fromIndex := pr.keyIndex(from)
	toIndex := pr.keyIndex(to)
	pr.updateMatrix(fromIndex, toIndex)
	pr.updateLinksNum(fromIndex)
}

// Rank calculate rank
func (pr *Pagerank) Rank(followingProb, tolerance float64, resultFunc func(label string, rank float64)) {
	size := pr.Len()
	invSize := 1.0 / float64(size)
	overSize := (1.0 - followingProb) / float64(size)
	existLinks := pr.getExistLinks()

	result := make([]float64, 0, size)
	for i := 0; i != size; i++ {
		result = append(result, invSize)
	}

	for change := 1.0; change >= tolerance; {
		newResult := pr.step(followingProb, overSize, existLinks, result)
		change = calculateTolerance(result, newResult)
		result = newResult
	}

	for i, v := range result {
		resultFunc(pr.Rekeymap[uint64(i)], v)
	}
}

// keyIndex
func (pr *Pagerank) keyIndex(key string) uint64 {
	index, ok := pr.Keymap[key]

	if !ok {
		index = uint64(len(pr.Rekeymap))
		pr.Rekeymap = append(pr.Rekeymap, key)
		pr.Keymap[key] = index
	}

	return index
}

// updateMatrix
func (pr *Pagerank) updateMatrix(fromAsIndex, toAsIndex uint64) {
	if missingSlots := len(pr.Keymap) - len(pr.Matrix); missingSlots > 0 {
		pr.Matrix = append(pr.Matrix, make([][]uint64, missingSlots)...)
	}
	pr.Matrix[toAsIndex] = append(pr.Matrix[toAsIndex], fromAsIndex)
}

// updateLinksNum
func (pr *Pagerank) updateLinksNum(fromAsIndex uint64) {
	if missingSlots := len(pr.Keymap) - len(pr.Links); missingSlots > 0 {
		pr.Links = append(pr.Links, make([]uint64, missingSlots)...)
	}
	pr.Links[fromAsIndex] += 1
}

// getExistLinks
func (pr *Pagerank) getExistLinks() []int {
	danglingNodes := make([]int, 0, len(pr.Links))

	for i, numberOutLinksForI := range pr.Links {
		if numberOutLinksForI == 0 {
			danglingNodes = append(danglingNodes, i)
		}
	}
	return danglingNodes
}

// step
func (pr *Pagerank) step(followingProb, overSize float64, existLinks []int, result []float64) []float64 {
	sumLinks := 0.0
	for _, v := range existLinks {
		sumLinks += result[v]
	}
	sumLinks /= float64(len(result))

	vsum := 0.0
	newResult := make([]float64, len(result))

	for i, from := range pr.Matrix {
		ksum := 0.0

		for _, index := range from {
			ksum += result[index] / float64(pr.Links[index])
		}

		newResult[i] = followingProb*(ksum+sumLinks) + overSize
		vsum += newResult[i]
	}

	for i := range newResult {
		newResult[i] /= vsum
	}

	return newResult
}

// calculateTolerance
func calculateTolerance(result, newResult []float64) float64 {
	acc := 0.0
	for i, v := range result {
		acc += math.Abs(v - newResult[i])
	}
	return acc
}
