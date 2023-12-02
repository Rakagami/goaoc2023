package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// This is super inefficient and a bit ugly, but I don't care at this point
// Returns the digit that it found, the length of the string and an error
// as a tuple. The prettiest(?) way to implement this is probably a Non-
// Deterministic Finite Automaton.
func checkWrittenDigit(s string) (int, int, error) {
    digits := []string {
        "zero",
        "one",
        "two",
        "three",
        "four",
        "five",
        "six",
        "seven",
        "eight",
        "nine",
    }

    for i, sref := range digits {
        if strings.Contains(s, sref) {
            return i, len(sref), nil
        }
    }

    return -1, -1, errors.New("Nothing found")
}

func parseLine(line string) (int, error) {
    digits := "0123456789"
    acc := ""
    unmatched_first := 0 // First inclusive index of unmatched string
    unmatched_last := 0 // Last exclusinve index of unmatched string
    for i, c := range line {
        if strings.Contains(digits, string(c)) {
            acc = acc + string(c)
            unmatched_first = i + 1
            unmatched_last = i + 1
        } else {
            unmatched_last += 1
            i, l, err := checkWrittenDigit(line[unmatched_first:unmatched_last])
            if err == nil {
                acc = acc + strconv.Itoa(i)
                unmatched_first = unmatched_last - l + 1
            }
        }
    }
    if len(acc) < 1 {
        return -1, errors.New("Not enough digits")
    } else {
        i, err := strconv.Atoi(string(acc[0]) + string(acc[len(acc)-1]))
        return i, err
    }

}

func main() {
    filepath := "./input.txt"
    f, _ := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
    defer f.Close()

    sc := bufio.NewScanner(f)
    sum := 0
    for sc.Scan() {
        s := sc.Text()
        i, err := parseLine(s)
        // fmt.Printf("Input string: %s; output: i: %d, err: %v\n", s, i, err)
        if err == nil {
            sum += i
        }
    }
    fmt.Printf("Final sum: %d\n", sum)
}
