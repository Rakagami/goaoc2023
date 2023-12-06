package utils

// IntSet

type IntSet struct {
    set map[int]bool
}

func NewIntSet() IntSet {
    intSet := IntSet{}
    intSet.set = map[int]bool{}
    return intSet
}

func (intSet IntSet) IntSetAdd(i int) {
    intSet.set[i] = true
}

func (intSet IntSet) IntSetRemove(i int) {
    delete(intSet.set, i)
}

// Returns union between intSet1 and intSet2
func (intSet1 IntSet) Union(intSet2 IntSet) IntSet{
    union := map[int]bool{}
    for k := range intSet1.set {
        union[k] = true
    }
    for k := range intSet2.set {
        union[k] = true
    }
    return IntSet{union}
}

// Returns intersection between intSet1 and intSet2
func (intSet1 IntSet) Intersect(intSet2 IntSet) IntSet{
    if len(intSet1.set) > len(intSet2.set) {
        intSet1, intSet2 = intSet2, intSet1
    }
    intersection := map[int]bool{}
    for k := range intSet1.set {
        if intSet2.set[k] {
            intersection[k] = true
        }
    }
    return IntSet{intersection}
}

// Returns true if intSet1 is subset of intSet2
func (intSet1 IntSet) IsSubset(intSet2 IntSet) bool{
    for k := range intSet1.set {
        if !intSet2.set[k] {
            return false
        }
    }
    return true
}
