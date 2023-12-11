package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/Rakagami/goaoc2023/utils"
)

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
func FindFittingMappedRange(ref int, ranges []MapRange) (utils.Range, error) {
    nextLowestRange, err := FindNextLowest(ref, ranges)
    if err != nil {
        return utils.Range{}, err
    } else if nextLowestRange.src + nextLowestRange.n <= ref {
        return utils.Range{}, err
    } else {
        return utils.Range{
            Start: nextLowestRange.dst,
            N: ref - nextLowestRange.src + 1,
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
func MapsToRange(in utils.Range, mapRanges []MapRange) utils.RangeSet {
    //fmt.Printf("--Maps To Range--\n")
    //fmt.Printf("\t[in]: %v; [mr]: %v\n", in, mapRanges)
    out := utils.RangeSet{}

    // inclusive indices
    start, end := in.Start, in.Start + in.N - 1
    cur := end

    for cur >= start {
        nextLowestMapRange, err := FindNextLowest(cur, mapRanges)

        // On no intersection, just add identity range
        if err != nil {
            out = out.Add(utils.Range{
                Start: start,
                N: cur - start + 1,
            })
            //fmt.Printf("\tNo valid map, just adding identity...\n")
            break
        } else if nextLowestMapRange.src + nextLowestMapRange.n - 1 < start {
            out = out.Add(utils.Range{
                Start: start,
                N: cur - start + 1,
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
            out = out.Add(utils.Range{
                Start: nextLowestMapRange.dst + (start - src),
                N: cur - start + 1,
            })
            break
        } else if start < src && cur <= lastSrcIdx {
            // 2
            //fmt.Printf("\tCase 2\n")
            out = out.Add(utils.Range{
                Start: nextLowestMapRange.dst,
                N: cur - src + 1,
            })
            cur = src - 1
            //fmt.Printf("\tcurrent out: %v\n", out)
            continue
        } else if start < src && cur > lastSrcIdx {
            // 3
            //fmt.Printf("\tCase 3\n")
            lastDstIdx := nextLowestMapRange.dst + nextLowestMapRange.n - 1
            out = out.Add(utils.Range{
                Start: lastDstIdx + 1,
                N: cur - lastSrcIdx,
            })
            out = out.Add(utils.Range{
                Start: nextLowestMapRange.dst,
                N: nextLowestMapRange.n,
            })
            cur = src - 1
            //fmt.Printf("\tcurrent out: %v\n", out)
            continue
        } else if start >= src && cur > lastSrcIdx {
            // 4
            //fmt.Printf("\tCase 4\n")
            lastDstIdx := nextLowestMapRange.dst + nextLowestMapRange.n - 1
            out = out.Add(utils.Range{
                Start: lastDstIdx + 1,
                N: cur - lastSrcIdx,
            })
            out = out.Add(utils.Range{
                Start: nextLowestMapRange.dst + (start - src),
                N: nextLowestMapRange.n - (start - src),
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
func MapsToRangeSet(in utils.RangeSet, mapRanges []MapRange) utils.RangeSet {
    out := utils.RangeSet{}

    for _, rg := range in {
        //fmt.Printf("\t\tBefore Union: %v\n", out)
        out = out.Union(MapsToRange(rg, mapRanges))
        //fmt.Printf("\t\tAfter Union: %v\n", out)
    }

    return out
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
            nums, _ := utils.ParseInts(s)
            seeds = nums
            break
        default:
            if strings.Contains(s, ":") {
                break
            } else {
                nums, _ := utils.ParseInts(s)
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

    seedRanges := utils.RangeSet{}
    for i := 0; i < len(seeds) / 2; i++ {
        seedRanges = append(seedRanges, utils.Range{
            Start: seeds[2*i],
            N: seeds[2*i+1],
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
    
    fmt.Printf("Minimum Result: %d\n", in[0].Start)
}
