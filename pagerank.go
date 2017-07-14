package pagerank

import (
	"math"
)

type Pagerank struct {
	Matrix   [][]uint64        // 节点被哪些节点指向
	Links    []uint64          // 节点指向数量
	Keymap   map[string]uint64 // 外键映射索引
	Rekeymap []string          // 索引映射外键
}

func NewPagerank() *Pagerank {
	pr := &Pagerank{}
	pr.Matrix = [][]uint64{}
	pr.Links = []uint64{}
	pr.Keymap = map[string]uint64{}
	pr.Rekeymap = []string{}
	return pr
}

// 页面总数
func (pr *Pagerank) Len() int {
	return len(pr.Rekeymap)
}

// 创建连接指向
func (pr *Pagerank) Link(from, to string) {
	fromIndex := pr.keyIndex(from)
	toIndex := pr.keyIndex(to)
	pr.updateMatrix(fromIndex, toIndex)
	pr.updateLinksNum(fromIndex)
}

// 开始计算分数
func (pr *Pagerank) Rank(followingProb, tolerance float64, resultFunc func(label string, rank float64)) {
	size := len(pr.Rekeymap)
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

// 键到内部索引绑定
func (pr *Pagerank) keyIndex(key string) uint64 {
	index, ok := pr.Keymap[key]

	if !ok {
		index = uint64(len(pr.Rekeymap))
		pr.Rekeymap = append(pr.Rekeymap, key)
		pr.Keymap[key] = index
	}

	return index
}

// 更新矩阵
func (pr *Pagerank) updateMatrix(fromAsIndex, toAsIndex uint64) {
	if missingSlots := len(pr.Keymap) - len(pr.Matrix); missingSlots > 0 {
		pr.Matrix = append(pr.Matrix, make([][]uint64, missingSlots)...)
	}
	pr.Matrix[toAsIndex] = append(pr.Matrix[toAsIndex], fromAsIndex)
}

// 更新节点指向数量
func (pr *Pagerank) updateLinksNum(fromAsIndex uint64) {
	if missingSlots := len(pr.Keymap) - len(pr.Links); missingSlots > 0 {
		pr.Links = append(pr.Links, make([]uint64, missingSlots)...)
	}
	pr.Links[fromAsIndex] += 1
}

// 获取所有 有指向的节点
func (pr *Pagerank) getExistLinks() []int {
	danglingNodes := make([]int, 0, len(pr.Links))

	for i, numberOutLinksForI := range pr.Links {
		if numberOutLinksForI == 0 {
			danglingNodes = append(danglingNodes, i)
		}
	}
	return danglingNodes
}

// 计算步骤
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

// 计算公差
func calculateTolerance(result, newResult []float64) float64 {
	acc := 0.0
	for i, v := range result {
		acc += math.Abs(v - newResult[i])
	}
	return acc
}
