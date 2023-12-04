package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Enum Pattern
type Symbol int
const (
    Digit Symbol = iota + 1
    Special
    Ignore
)

type Coordinate struct {
    x, y int
}

type ValuedCoordinate struct {
    x, y, i, n int
}

//Parses a symbolMap and returns a list of coordinates of special symbols
func sparseMatrixParser(symbolMap [][]byte, width int, height int) []Coordinate{
    digits := "0123456789"
    coordinates := []Coordinate{}
    for i := 0; i < height; i++ {
        for j := 0; j < width; j++ {
            if strings.Contains(digits, string(symbolMap[i][j])) {
                continue
            } else if symbolMap[i][j] == '.' {
                continue
            } else {
                coordinates = append(coordinates, Coordinate{j, i})
            }
        }
    }
    return coordinates
}

// Parses a number in the string
func parseNumber(symbolMap [][]byte, coordinate Coordinate) (ValuedCoordinate, error) {
    digits := "0123456789"
    line := symbolMap[coordinate.y]
    lineLength := len(line)

    leftMost := coordinate.x
    rightMost := coordinate.x

    isDigit := func(b byte) bool {
        return strings.Contains(digits, string(b))
    }

    if !isDigit(line[leftMost]) {
        return ValuedCoordinate{}, errors.New("Not Digit")
    } else if coordinate.x < 0 || coordinate.x >= lineLength || coordinate.y < 0 || coordinate.y >= len(symbolMap) {
        return ValuedCoordinate{}, errors.New("Not Valid Coordinate")
    }

    // Go left until the leftMost digit coordinate is found
    for ;leftMost >= 0; leftMost--{
        if leftMost == 0 {
            break
        } else if !isDigit(line[leftMost-1]) {
            break
        }
    }

    // Go right until the rightMost digit coordinate is found
    for ; rightMost < lineLength; rightMost++ {
        if rightMost == lineLength - 1 {
            break
        } else if !isDigit(line[rightMost+1]) {
            break
        }
    }

    numStr := line[leftMost:rightMost + 1]
    num, err := strconv.Atoi(string(numStr))

    return ValuedCoordinate{
        x: rightMost,
        y: coordinate.y,
        i: num,
        n: len(numStr),
    }, err
}

//Check all 8 possible adjacent numbers for values and returns the set of ValuedCoordinates
func findAdjacentNumbers(symbolMap [][]byte, coordinate Coordinate, coordinateSet map[ValuedCoordinate]bool) map[ValuedCoordinate]bool{
    validityCheck := func(c Coordinate) bool {
        if c.x < 0 || c.y < 0 || (c.x == coordinate.x && c.y == coordinate.y) {
            return false
        } else {
            return true
        }
    }
    for i := coordinate.y - 1; i <= coordinate.y + 1; i++ {
        for j := coordinate.x - 1; j <= coordinate.x + 1; j++ {
            if !validityCheck(Coordinate{j, i}) {
                continue
            }

            vc, err := parseNumber(symbolMap, Coordinate{j, i})
            if err != nil {
                continue
            } else {
                coordinateSet[vc] = true
            }
        }
    }
    return coordinateSet
}

// We'll just assume that the width is always the same
func readFileToMat(path string) ([][]byte, int, int, error) {
    f, _ := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
    defer f.Close()
    sc := bufio.NewScanner(f)
    symbolMap := [][]byte{}
    width := 0
    height := 0
    for sc.Scan() {
        s := sc.Text()
        if width == 0 && width != len(s) {
            width = len(s)
        } else if width != len(s) {
            return [][]byte{}, -1, -1, errors.New("Invalid input file")
        }

        symbolMap = append(symbolMap, []byte(s))
        height++
    }
    return symbolMap, width, height, nil
}

func computeGearRatio(symbolMap [][]byte, coordinate Coordinate) (int, error) {
    if symbolMap[coordinate.y][coordinate.x] != byte('*') {
        return 0, errors.New("Not a gear symbol")
    }
    coordinateSet := map[ValuedCoordinate]bool{}
    coordinateSet = findAdjacentNumbers(symbolMap, coordinate, coordinateSet)
    if len(coordinateSet) != 2 {
        return 0, errors.New("Not enough part numbers")
    }

    product := 1
    for vc := range coordinateSet {
        product *= vc.i
    }

    return product, nil
}

func main() {
    filepath := "./input.txt"

    symbolMap, width, height, err := readFileToMat(filepath)
    if err != nil {
        fmt.Printf("For some reason, an error has occured: %v", err)
    }
    //fmt.Printf("Filemap: %v\n", symbolMap)

    specialCoordinates := sparseMatrixParser(symbolMap, width, height)
    coordinateSet := map[ValuedCoordinate]bool{}

    gSum := 0
    for _, sc := range specialCoordinates {
        coordinateSet = findAdjacentNumbers(symbolMap, sc, coordinateSet)
        gr, _ := computeGearRatio(symbolMap, sc)
        gSum += gr
    }

    //fmt.Println(coordinateSet)
    sum := 0
    for vc := range coordinateSet {
        fmt.Println(vc)
        sum += vc.i
    }
    
    fmt.Printf("Total Sum of Part Numbers: %d\n", sum)
    fmt.Printf("Total Sum of Gear Ratio: %d\n", gSum)

    //fmt.Printf("%v\n", specialCoordinates)
}
