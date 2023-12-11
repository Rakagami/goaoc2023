package main

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Rakagami/goaoc2023/utils"
)

type Range struct {
    start int // inclusive start
    n int // length of range
}
// Last value of range
func (rs Range) End() int {return rs.start + rs.n - 1}
// Checking whether value is in range
func (rs Range) Check(i int) bool {return i >= rs.start && i < rs.start + rs.n}

type RangeSet []Range
func (rs RangeSet) Len() int      {return len(rs)}
func (rs RangeSet) Swap(i, j int) {rs[i], rs[j] = rs[j], rs[i]}
func (rs RangeSet) Less(i, j int) bool {return rs[i].start < rs[j].start}
func (rs RangeSet) Sort() {
    sort.Sort(rs)
}

// Linear Search because I'm lazy. Returns index
func (rs RangeSet) FindRange(ref int) (int, error) {
    for i, rg := range rs {
        if ref <= rg.start && ref < rg.start + rg.n {
            return i, nil
        }
    }
    return -1, errors.New("Not Found")
}

// This function assumes that the rs RangeSet is already sorted
// Uses binary search
func (rs RangeSet) FindNextLowest(ref int) (int, error) {
    if rs.Len() == 0 || (rs.Len() == 1 && rs[0].start > ref) {
        return -1, errors.New("Not found")
    } else if rs.Len() == 1 && rs[0].start <= ref {
        return 0, nil
    }

    pivotIdx := int(rs.Len() / 2)
    if rs[pivotIdx].start <= ref {
        idx, err := rs[pivotIdx:].FindNextLowest(ref)
        return pivotIdx + idx, err
    } else {
        return rs[:pivotIdx].FindNextLowest(ref)
    }
}

// A SoftIntersection returns the a RangeSet of ranges that fall into the range
// This function assumes that the rs RangeSet is already sorted
func (rs RangeSet) SoftIntersection(rg Range) (int, int, error) {
    start, end := rg.start, rg.End()

    lastIdx, err := rs.FindNextLowest(end)
    firstIdx := lastIdx

    // Cases if there is no overlap
    if err != nil {
        return -1, -1, errors.New("No intersection")
    } else if rs[lastIdx].End() < start {
        return -1, -1, errors.New("No intersection")
    }

    // Finding the range in which start fits
    for rs[firstIdx].start > start && firstIdx > 0 {
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
        startVal, endVal = rg.start, rg.End()
    } else if rg.start < rs[startRgIdx].start {
        startVal, endVal = rg.start, rs[endRgIdx].End()
        rs = rs.RemoveIdxRange(startRgIdx, endRgIdx)
    } else if rg.End() > rs[endRgIdx].End() {
        startVal, endVal = rs[startRgIdx].start, rg.End()
        rs = rs.RemoveIdxRange(startRgIdx, endRgIdx)
    } else {
        startVal, endVal = rs[startRgIdx].start, rs[endRgIdx].End()
        rs = rs.RemoveIdxRange(startRgIdx, endRgIdx)
    }

    rs = append(rs, Range{
        start: startVal,
        n: endVal - startVal + 1,
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

type MapRange struct {
    src int
    dst int
    n int
}

type MapRanges []MapRange
func (mrs MapRanges) Len() int      {return len(mrs)}
func (mrs MapRanges) Swap(i, j int) {mrs[i], mrs[j] = mrs[j], mrs[i]}
func (mrs MapRanges) Less(i, j int) bool {return mrs[i].src < mrs[j].src}

// Uses binary search to find next lowest integer with depth
// Unbounded recursive function, let's goooooo
func FindNextLowest(ref int, ranges []MapRange) (MapRange, error) {
    if len(ranges) == 0 || (len(ranges) == 1 && ranges[0].src > ref) {
        return MapRange{}, errors.New("Not found")
    } else if len(ranges) == 1 && ranges[0].src <= ref {
        return ranges[0], nil
    }

    pivotIdx := int(len(ranges) / 2)
    if ranges[pivotIdx].src <= ref {
        return FindNextLowest(ref, ranges[pivotIdx:])
    } else {
        return FindNextLowest(ref, ranges[:pivotIdx])
    }
}

// Given reference index ref, this function returns the range that is smaller
// than the reference. The returned range is in the destination space.
func FindFittingMappedRange(ref int, ranges []MapRange) (Range, error) {
    nextLowestRange, err := FindNextLowest(ref, ranges)
    if err != nil {
        return Range{}, err
    } else if nextLowestRange.src + nextLowestRange.n <= ref {
        return Range{}, err
    } else {
        return Range{
            start: nextLowestRange.dst,
            n: ref - nextLowestRange.src + 1,
        }, nil
    }
}

// Mapping ranges, assumes that ranges is already sorted
func MapsToSingle(in int, ranges []MapRange) int {
    nextLowestRange, err := FindNextLowest(in, ranges)
    //fmt.Printf("\t\tfoudn next lowest: %v\n", nextLowestRange)
    if err != nil {
        return in
    } else if nextLowestRange.src + nextLowestRange.n <= in {
        return in
    } else {
        return nextLowestRange.dst + (in - nextLowestRange.src)
    }
}

// Maps range to set of ranges. Iirc there can be no overlap in the image of a
// map, so that's something we don't have to worry about
func MapsToRange(in Range, mapRanges []MapRange) RangeSet {
    //fmt.Printf("--Maps To Range--\n")
    //fmt.Printf("\t[in]: %v; [mr]: %v\n", in, mapRanges)
    out := RangeSet{}

    // inclusive indices
    start, end := in.start, in.start + in.n - 1
    cur := end

    for cur >= start {
        nextLowestMapRange, err := FindNextLowest(cur, mapRanges)

        // On no intersection, just add identity range
        if err != nil {
            out = out.Add(Range{
                start: start,
                n: cur - start + 1,
            })
            //fmt.Printf("\tNo valid map, just adding identity...\n")
            break
        } else if nextLowestMapRange.src + nextLowestMapRange.n - 1 < start {
            out = out.Add(Range{
                start: start,
                n: cur - start + 1,
            })
            //fmt.Printf("\tNo valid map, just adding identity...\n")
            break
        } 

        //fmt.Printf("\tNext lowest map range:%v\n", nextLowestMapRange)

        lastSrcIdx := nextLowestMapRange.src + nextLowestMapRange.n - 1
        src := nextLowestMapRange.src

        // There are four overlap scenarios:
        //
        // 1.
        //                    Start            Cur
        //                     v                v
        //    [in]  :          ##################
        //               Src                         lastSrcIdx
        //                v                             v
        //    [map] :     ###############################
        //
        // 2.
        //                    Start            Cur
        //                     v                v
        //    [in]  :          ##################
        //                          Src         lastSrcIdx
        //                           v              v
        //    [map] :                ################
        //
        // 3.
        //                    Start            Cur
        //                     v                v
        //    [in]  :          ##################
        //                        Src     lastSrcIdx
        //                         v        v
        //    [map] :              ##########
        //
        // 4.
        //                    Start            Cur
        //                     v                v
        //    [in]  :          ##################
        //               Src          lastSrcIdx
        //                v             v
        //    [map] :     ###############

        if start >= src && cur <= lastSrcIdx {
            // 1
            //fmt.Printf("\tCase 1\n")
            out = out.Add(Range{
                start: nextLowestMapRange.dst + (start - src),
                n: cur - start + 1,
            })
            break
        } else if start < src && cur <= lastSrcIdx {
            // 2
            //fmt.Printf("\tCase 2\n")
            out = out.Add(Range{
                start: nextLowestMapRange.dst,
                n: cur - src + 1,
            })
            cur = src - 1
            //fmt.Printf("\tcurrent out: %v\n", out)
            continue
        } else if start < src && cur > lastSrcIdx {
            // 3
            //fmt.Printf("\tCase 3\n")
            lastDstIdx := nextLowestMapRange.dst + nextLowestMapRange.n - 1
            out = out.Add(Range{
                start: lastDstIdx + 1,
                n: cur - lastSrcIdx,
            })
            out = out.Add(Range{
                start: nextLowestMapRange.dst,
                n: nextLowestMapRange.n,
            })
            cur = src - 1
            //fmt.Printf("\tcurrent out: %v\n", out)
            continue
        } else if start >= src && cur > lastSrcIdx {
            // 4
            //fmt.Printf("\tCase 4\n")
            lastDstIdx := nextLowestMapRange.dst + nextLowestMapRange.n - 1
            out = out.Add(Range{
                start: lastDstIdx + 1,
                n: cur - lastSrcIdx,
            })
            out = out.Add(Range{
                start: nextLowestMapRange.dst + (start - src),
                n: nextLowestMapRange.n - (start - src),
            })
            break
        } else {
            // This case should not occur
            panic(errors.New("Unexpected state reached"))
        }
    }

    //fmt.Printf("\tout: %v\n", out)

    //fmt.Printf("--Maps To End--\n")

    return out
}

// Mapping ranges, assumes that mapRanges is already sorted
func MapsToRangeSet(in RangeSet, mapRanges []MapRange) RangeSet {
    out := RangeSet{}

    for _, rg := range in {
        //fmt.Printf("\t\tBefore Union: %v\n", out)
        out = out.Union(MapsToRange(rg, mapRanges))
        //fmt.Printf("\t\tAfter Union: %v\n", out)
    }

    return out
}

func parseInts(s string) ([]int, error) {
    re := regexp.MustCompile("\\d+")
    
    matchArr := re.FindAll([]byte(s), -1)
    intArr := []int{}
    for _, match := range matchArr {
        i, err := strconv.Atoi(string(match))
        if err != nil {
            continue
        }
        intArr = append(intArr, i)
    }

    return intArr, nil
}

const TOTAL_MAPS int = 7

func main() {
    seeds := []int{}
    rangesList := [TOTAL_MAPS]MapRanges{}
    parseState := 0
    parseLine := func(s string) error {
        if s == "" {
            parseState++
            return nil
        }
        
        switch parseState {
        case 0:
            nums, _ := parseInts(s)
            seeds = nums
            break
        default:
            if strings.Contains(s, ":") {
                break
            } else {
                nums, _ := parseInts(s)
                mapRange := MapRange{
                    src: nums[1],
                    dst: nums[0],
                    n: nums[2],
                }
                rangesList[parseState - 1] = append(rangesList[parseState - 1], mapRange)
            }
        }

        return nil
    }

    utils.IterateLines("input.txt", parseLine)

    seedRanges := RangeSet{}
    for i := 0; i < len(seeds) / 2; i++ {
        seedRanges = append(seedRanges, Range{
            start: seeds[2*i],
            n: seeds[2*i+1],
        })
    }

    // Sort MapRanges
    for i := 0; i < len(rangesList); i++ {
        sort.Sort(rangesList[i])
    }

    fmt.Printf("Seed Ranges: %v\n", seedRanges)
    fmt.Printf("Sorted Map Ranges: %v\n", rangesList)

    in := seedRanges
    for i:=0; i<TOTAL_MAPS; i++ {
        fmt.Printf("Iteration %d; [in]: %v\n", i, in)
        //fmt.Printf("\tranges: %v\n", rangesList[i])
        in = MapsToRangeSet(in, rangesList[i])
        //fmt.Printf("\tin: %d\n", in)
    }

    in.Sort()

    fmt.Printf("Final [in]: %v\n", in)

    fmt.Printf("\n")
    
    fmt.Printf("Minimum Result: %d\n", in[0].start)
}
