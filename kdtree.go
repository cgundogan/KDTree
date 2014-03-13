package kdtree

import (
	"sort"
        "math"
        "fmt"
)

/* the tree is structured by node pointers */
type Node struct {
	dim    int
	median Vector
	left   *Node
	right  *Node
}

type KDTree struct {
	root *Node
}

func NewNode(dim int, median Vector, left, right *Node) *Node {
	return &Node{dim, median, left, right}
}

/* the sorted tree is recursively built */
func build(dataSet DataSet) *Node {
	if len(dataSet.set) == 1 {
		return NewNode(dataSet.dim, dataSet.set[0], nil, nil)
	}

	sort.Sort(ByDim(dataSet))
	left, right, median := dataSet.splitByMedianAndDim()
        
        var leftNode, rightNode *Node = nil, nil
        leftNode = build(left)
        if len(right.set) != 0 {
	        rightNode = build(right)
        }
	
        return NewNode(dataSet.dim, median, leftNode, rightNode)
}

func NewKDTree(path string) (KDTree, error) {
	dataSet, err := Load(path)
        if err != nil {
                return KDTree{}, err
        }
        return NewKDTreeByDataSet(dataSet), err
}

func NewKDTreeByDataSet(dataSet DataSet) (KDTree) {
        kdTree := KDTree{}
        kdTree.root = build(dataSet)
        return kdTree
}

/* this function calculates the euclidean distance */
func calcDist(vec1, vec2 Vector) float64 {
        dist := 0.0
        for i := range vec1 {
                dist += math.Pow((vec1[i] - vec2[i]), 2)
        }
        return dist
}

/* the tree will be traversed by comparing the values until a leaf is found.
 * afterwards the distances will be calculated and compared and the node
 * with the lowest distance is the nearest vector
 */
func find(node *Node, vector Vector) *Node {
        var bestNode *Node = nil
        isLeft := false
        
        if node == nil {
                return node
        } else if vector[node.dim] <= node.median[node.dim] {
                bestNode = find(node.left, vector)
                isLeft = true
        } else {
                bestNode = find(node.right, vector)
        }

        if bestNode == nil || calcDist(node.median, vector) < calcDist(bestNode.median, vector) {
                bestNode = node
        }
        /* the sibling might be checked if the current best node is near the median */
        if calcDist(Vector{node.median[node.dim]}, Vector{vector[node.dim]}) < calcDist(bestNode.median, vector) {
                var siblingNode *Node = nil
                
                if isLeft {
                        siblingNode = find(node.right, vector)
                } else {
                        siblingNode = find(node.left, vector)
                }
                
                if siblingNode != nil && calcDist(siblingNode.median, vector) < calcDist(bestNode.median, vector) {
                        bestNode = siblingNode
                }
        }
        
        return bestNode
}

/* finds the nearest node to the given vector */
func (self *KDTree) FindNearest(vector Vector) Vector {
        nearestNode := find(self.root, vector)
        if (nearestNode != nil) {
                return nearestNode.median
        } else {
                return nil
        }
}

/* the printing is a little bit odd for unbalanced trees,.. */
func PrintTree(node *Node, prefix string, isLeft bool, isFirst bool) {
	pre := ""
	if isLeft {
		pre	= "    ┌────"
	} else {
		pre	= "    └────"
	}

	if isFirst {
		pre	= "    "
	}

	if (node.left != nil) {
		PrintTree(node.left, prefix +  "     ", true, false)
	}

	fmt.Printf("%v (%v, %v)\n", prefix + pre, node.dim, node.median)

	if (node.right != nil) {
		PrintTree(node.right, prefix + "     ", false, false)
	}
}
