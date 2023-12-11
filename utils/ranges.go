package utils

import (
    "sort"
    "errors"
)

type Range struct {
    Start int // inclusive start
    N int // length of range
}
// Last value of range
func (rs Range) End() int {return rs.Start + rs.N - 1}
// Checking whether value is in range
func (rs Range) Check(i int) bool {return i >= rs.Start && i < rs.Start + rs.N}

type RangeSet []Range

func (rs RangeSet) Len() int      {return len(rs)}
func (rs RangeSet) Swap(i, j int) {rs[i], rs[j] = rs[j], rs[i]}
func (rs RangeSet) Less(i, j int) bool {return rs[i].Start < rs[j].Start}
func (rs RangeSet) Sort() {
    sort.Sort(rs)
}

// Linear Search because I'm lazy. Returns index
func (rs RangeSet) FindRange(ref int) (int, error) {
    for i, rg := range rs {
        if ref <= rg.Start && ref < rg.Start + rg.N {
            return i, nil
        }
    }
    return -1, errors.New("Not Found")
}

// This function assumes that the rs RangeSet is already sorted
// Uses binary search
func (rs RangeSet) FindNextLowest(ref int) (int, error) {
    if rs.Len() == 0 || (rs.Len() == 1 && rs[0].Start > ref) {
        return -1, errors.New("Not found")
    } else if rs.Len() == 1 && rs[0].Start <= ref {
        return 0, nil
    }

    pivotIdx := int(rs.Len() / 2)
    if rs[pivotIdx].Start <= ref {
        idx, err := rs[pivotIdx:].FindNextLowest(ref)
        return pivotIdx + idx, err
    } else {
        return rs[:pivotIdx].FindNextLowest(ref)
    }
}

// A SoftIntersection returns the a RangeSet of ranges that fall into the range
// This function assumes that the rs RangeSet is already sorted
func (rs RangeSet) SoftIntersection(rg Range) (int, int, error) {
    start, end := rg.Start, rg.End()

    lastIdx, err := rs.FindNextLowest(end)
    firstIdx := lastIdx

    // Cases if there is no overlap
    if err != nil {
        return -1, -1, errors.New("No intersection")
    } else if rs[lastIdx].End() < start {
        return -1, -1, errors.New("No intersection")
    }

    // Finding the range in which start fits
    for rs[firstIdx].Start > start && firstIdx > 0 {
        firstIdx--
    }

    if rs[firstIdx].Check(start) || firstIdx == lastIdx {
        return firstIdx, lastIdx, nil
    } else {
        // Special handling in case the first range doesn't contain start
        return firstIdx + 1, lastIdx, nil
    }

}

// Removes indices from i to j
func (rs RangeSet) RemoveIdxRange(i int, j int) RangeSet {
    rs = append(rs[:i], rs[j+1:]...)
    return rs
}

// Adds a range to range set. Merges ranges if necessary
func (rs RangeSet) Add(rg Range) RangeSet {
    rs.Sort()
    startRgIdx, endRgIdx, err := rs.SoftIntersection(rg)

    startVal, endVal := -1, -1

    if err != nil {
        // If there's no overlap, just add range
        startVal, endVal = rg.Start, rg.End()
    } else if rg.Start < rs[startRgIdx].Start {
        startVal, endVal = rg.Start, rs[endRgIdx].End()
        rs = rs.RemoveIdxRange(startRgIdx, endRgIdx)
    } else if rg.End() > rs[endRgIdx].End() {
        startVal, endVal = rs[startRgIdx].Start, rg.End()
        rs = rs.RemoveIdxRange(startRgIdx, endRgIdx)
    } else {
        startVal, endVal = rs[startRgIdx].Start, rs[endRgIdx].End()
        rs = rs.RemoveIdxRange(startRgIdx, endRgIdx)
    }

    rs = append(rs, Range{
        Start: startVal,
        N: endVal - startVal + 1,
    })

    rs.Sort()
    return rs
}

func (rs1 RangeSet) Union(rs2 RangeSet) RangeSet {
    for _, rg := range rs2 {
        rs1 = rs1.Add(rg)
    }
    return rs1
}
