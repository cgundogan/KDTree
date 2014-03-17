package kdtree

import (
	"bufio"
	"os"
	"strconv"
	"strings"
        "fmt"
)

type Vector []float64

/* for every step down in the kd-tree we have to split the dataset and choose another dimension.
 * set is a slice of Vector, which means there will be only one copy of the underlying array
 * and multiple "views" on it
 */
type DataSet struct {
	set    []Vector
	maxDim int
	dim    int
}

func NewVector(dim int) Vector {
	return make(Vector, dim)
}

func NewDataSet(set []Vector, maxDim, dim int) DataSet {
	return DataSet{set, maxDim, dim}
}

/* this function will parse a file containing vectors of the form:
 * x1|x2|..|xn
 * y1|y2|..|yn
 * ..
 * into a new DataSet
 */
func Load(path string) (DataSet, error) {
	file, err := os.Open(path)
	dataSet := NewDataSet([]Vector{}, 0, 0)
	if err != nil {
		return dataSet, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	first := true
	for scanner.Scan() {
		vals := strings.Split(scanner.Text(), "|")
		
                if first {
			dataSet.maxDim = len(vals)
			first = false
		}
                
                if dataSet.maxDim != len(vals) {
                        return DataSet{}, fmt.Errorf("dimension of all vectors must be equal!")
                }
		
                vector := NewVector(len(vals))
		for k, v := range vals {
			vector[k], err = strconv.ParseFloat(v, 64)
		}
		dataSet.set = append(dataSet.set, vector)
	}
	return dataSet, nil /* no error */
}

/* this function will only work on a sorted slice */
func (self *DataSet) medianIdx() int {
	return len(self.set) / 2
}

/* the following three functions are part of the sorting interface and must be implemented
 * to be able to sort the DataSet.
 * The slice of Vector will be sorted on only one specific dimension.
 */
type ByDim DataSet

func (a ByDim) Len() int {
	return len(a.set)
}

func (a ByDim) Swap(i, j int) {
	a.set[i], a.set[j] = a.set[j], a.set[i]
}

func (a ByDim) Less(i, j int) bool {
	return a.set[i][a.dim] < a.set[j][a.dim]
}

/* the DataSet will be split in a specific dimension by the median into two sub slices
 * both slices and the median will be returned
 */
func (self *DataSet) splitByMedianAndDim() (DataSet, DataSet, Vector) {
	medianIdx := self.medianIdx()

        /* choose next dimension incrementally */
	nextDim := (self.dim + 1) % self.maxDim
        
        /* split the data by creating two new slices */
	leftData, rightData := self.set[:medianIdx], self.set[medianIdx+1:]
	median := self.set[medianIdx]

        /* if there are only two vectors, choose the first to be the left data and
         * the second vector to be the median, right data is empty*/
	if len(self.set) == 2 {
		leftData, rightData = []Vector{self.set[0]}, self.set[0:0]
		median = self.set[1]
	}

	return NewDataSet(leftData, self.maxDim, nextDim),
		NewDataSet(rightData, self.maxDim, nextDim),
		median
}
