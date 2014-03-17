package kdtree

import "testing"

func vectorEquals(vec1 Vector, vec2 Vector) bool {
        if len(vec1) != len(vec2) {
                return false
        }
        for i:= range vec1 {
                if vec1[i] != vec2[i] {
                        return false
                }
        }
        return true
}

func TestKDTree(t *testing.T) {
	var tests = []struct {
		vecs []Vector
                nearestTo, want Vector
	}{
		{
                        []Vector{
                                Vector{         1.0,    2.3,    9.87            },
                                Vector{         0.8,    10.98,  10.09           },
                                Vector{         100.9,  0.1,    1.1             },
                        },
                        Vector{         1.0,    0.9,    1.1     },
                        Vector{         1.0,    2.3,    9.87    },
                },
	}
	for _, c := range tests {
                dataSet := NewDataSet(c.vecs, len(c.vecs[0]), 0)
                kdTree := NewKDTreeByDataSet(dataSet)
                got := kdTree.FindNearest(c.nearestTo)
		if !vectorEquals(got, c.want) {
			t.Errorf("kdTree.FindNearest(%v) == %v, want %v", c.vecs, got, c.want)
		}
	}
}

