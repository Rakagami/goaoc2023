package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

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

func (intSet1 IntSet) Union(intSet2 IntSet) IntSet{
    union := map[int]bool{}
    for k, _:= range intSet1.set {
        union[k] = true
    }
    for k, _:= range intSet2.set {
        union[k] = true
    }
    return IntSet{union}
}

func (intSet1 IntSet) Intersect(intSet2 IntSet) IntSet{
    if len(intSet1.set) > len(intSet2.set) {
        intSet1, intSet2 = intSet2, intSet1
    }
    intersection := map[int]bool{}
    for k, _ := range intSet1.set {
        if intSet2.set[k] {
            intersection[k] = true
        }
    }
    return IntSet{intersection}
}

// Returns true if intSet1 is subset of intSet2
func (intSet1 IntSet) IsSubset(intSet2 IntSet) bool{
    for k, _ := range intSet1.set {
        if !intSet2.set[k] {
            return false
        }
    }
    return true
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

func parseIntSet(s string) (IntSet, error) {
    intSet := NewIntSet()
    intArr, err := parseInts(s)
    if err != nil {
        return IntSet{}, err
    }

    for _, i := range intArr {
        intSet.IntSetAdd(i)
    }

    return intSet, nil
}

func iterateLines(filePath string, callback func(string) error) error {
    f, _ := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
    defer f.Close()
    sc := bufio.NewScanner(f)
    for sc.Scan() {
        s := sc.Text()
        err := callback(s)
        if err != nil {
            return err
        }
    }
    return nil
}

// Returns a new countArr
func applyWin(index int, matchArr []int, countArr []int) []int {
    n := len(matchArr)
    if index >= n - 1 {
        return countArr
    } else if matchArr[index] == 0 {
        return countArr
    } else {
        for i := index + 1; i < index + matchArr[index] + 1; i++ {
            countArr[i] = countArr[i] + 1
            countArr = applyWin(i, matchArr, countArr)
        }
    }

    return countArr
}

func main() {
    sum := 0

    matchArr := []int{}
    countArr := []int{}
    callback := func(line string) error {
        // splitting for Game ID part and BallSet part
        v := strings.SplitN(line, ":", 2)
        if len(v) != 2 {
            return errors.New("Unknown error")
        }
        numStr := strings.SplitN(v[1], "|", 2)
        winningStr, drawStr := numStr[0], numStr[1]
        winningSet, _ := parseIntSet(string(winningStr))
        drawSet, _ := parseIntSet(string(drawStr))

        matchingSet := winningSet.Intersect(drawSet)

        matchArr = append(matchArr, len(matchingSet.set))
        countArr = append(countArr, 1)
        if len(matchingSet.set) > 0 {
            sum += 1 << (len(matchingSet.set) - 1)
        }
        return nil
    }
    err := iterateLines("./input.txt", callback)

    //fmt.Println(matchArr)
    for i := len(matchArr) - 1; i >= 0; i-- {
        countArr = applyWin(i, matchArr, countArr)
        //fmt.Println(countArr)
    }

    countSum := 0
    for _, count := range countArr {
        countSum += count
    }

    fmt.Printf("Final Count Sum: %d\n", countSum)
    fmt.Printf("Final Winning Point Sum: %d\n", sum)
    fmt.Printf("Error of iterate lines: %v\n", err)
}
