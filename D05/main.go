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

// TODO
// Mapping ranges, assumes that mapRanges is already sorted
func MapsToRange(in Range, mapRanges []MapRange) []Range {
    retRanges := []Range{}
    idxLastRange := len(mapRanges) - 1
    idxFirstRange := 0

    findMapRange := func(mr MapRange) int {
        for i, v := range mapRanges {
            if v == mr {
                return i
            }
        }
        return -1
    }

    nextLowestOfEnd, err1 := FindNextLowest(in.start + in.n - 1, mapRanges)
    if err1 != nil {
        // The entire range is outside of the map ranges
        return append(retRanges, in)
    } else {
        // Setting the idx of the last overlapping range
        idxLastRange = findMapRange(nextLowestOfEnd)
    }

    nextLowestOfStart, err2 := FindNextLowest(in.start, mapRanges)
    if err2 == nil {
        // If there is a next lowest
        idxFirstRange = findMapRange(nextLowestOfStart)
    }

    for idx:=idxFirstRange; idx<idxLastRange+1; idx++ {
        out := Range{
            start: mapRanges[idx].dst,
            n: mapRanges[idx].n,
        }
        // Order of these if statements is important

        // Special Handling, input range start after map range starts
        if idx == idxFirstRange && err2 != nil {
            if in.start > nextLowestOfStart.src {
                out.start += in.start - nextLowestOfStart.src
                out.n -= in.start - nextLowestOfStart.src
            }
        }

        // Special Handling, if input range stops before map range stops
        if idx == idxLastRange {
            if in.start + in.n < mapRanges[idx].src + mapRanges[idx].n  {
                // Subtracting the difference
                out.n -= (mapRanges[idx].src + mapRanges[idx].n) - (in.start + in.n)
            }
        }

        retRanges = append(retRanges, out)
    }
    return retRanges
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

    // Sort MapRanges
    for i := 0; i < len(rangesList); i++ {
        sort.Sort(rangesList[i])
    }

    min_result := -1 // negative means not defined yet
    for _, seed := range seeds {
        in := seed
        fmt.Printf("seed: %d\n", seed)
        for i:=0; i<TOTAL_MAPS; i++ {
            //fmt.Printf("\tranges: %v\n", rangesList[i])
            in = MapsToSingle(in, rangesList[i])
            //fmt.Printf("\tin: %d\n", in)
        }
        fmt.Printf("\tfinal: %d\n", in)
        if min_result < 0 || in < min_result {
            min_result = in
        }
    }

    fmt.Printf("\n")
    
    fmt.Printf("Minimum Result: %d\n", min_result)
}
