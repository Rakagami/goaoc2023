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

type BallSet struct {
    n_red int
    n_blue int
    n_green int
}

// There are three important delimiters:
// : delimits Game ID from draws
// ; delimits draws
// , delimits types of balls within a draw

func parseInt(s string) (int, error) {
    re := regexp.MustCompile("\\d+")
    
    i_str := string(re.Find([]byte(s)))
    i, err := strconv.Atoi(i_str)
    if err != nil {
        return -1, errors.New("Not an int found")
    }
    return i, nil
}

func parseSetString(s string) (BallSet, error) {
    n_red := 0
    n_blue := 0
    n_green := 0

    var err error = nil

    re := regexp.MustCompile("\\d+ red|\\d+ blue|\\d+ green")
    match_list := re.FindAll([]byte(s), -1)

    for _, match := range(match_list) {
        s = string(match)
        if strings.Contains(s, "red") {
            n_red, err = parseInt(s)
        } else if strings.Contains(s, "green") {
            n_green, err = parseInt(s)
        } else if strings.Contains(s, "blue") {
            n_blue, err = parseInt(s)
        }
    }

    if err != nil {
        fmt.Printf("error %v\n", err)
        return BallSet{}, err
    }

    return BallSet{n_red, n_blue, n_green}, nil
}

func parseSetsString(s string) ([]BallSet, error) {
    // splitting for Game ID part and BallSet part
    v := strings.SplitN(s, ";", -1)

    sets := []BallSet{}

    for _, s := range v {
        b, err := parseSetString(s)
        if err != nil {
            continue
        }
        sets = append(sets, b)
    }

    return sets, nil
}

// Parses line
// Returns game id, list of ballsets, and error
func parseLine(s string) (int, []BallSet, error) {
    // splitting for Game ID part and BallSet part
    v := strings.SplitN(s, ":", 2)
    if len(v) != 2 {
        return -1, []BallSet{}, errors.New("Invalid Line")
    }
    
    id_string, set_string := v[0], v[1]

    id, err := parseInt(id_string)
    if err != nil {
        return -1, []BallSet{}, err
    }

    setlist, err := parseSetsString(set_string)
    if err != nil {
        return -1, []BallSet{}, err
    }

    fmt.Printf("\tid: %d; %v;\n", id, setlist)

    return id, setlist, nil
}

// If returns True, then b1 can be a subset of b2, False otherwise
func isSubset(b1 BallSet, b2 BallSet) bool {
    return b1.n_red <= b2.n_red && b1.n_blue <= b2.n_blue && b1.n_green <= b2.n_green
}

// Returns ID of the game if it's a possible game
func getValidGameID(s string, b_ref BallSet) (int, error) {
    id, setlist, err := parseLine(s)

    if err != nil {
        return -1, err
    }

    for _, b := range setlist {
        if !(isSubset(b, b_ref)) {
            return 0, nil
        }
    }

    return id, nil
}

func main() {
    filepath := "./input.txt"
    f, _ := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
    defer f.Close()

    b := BallSet{12, 14, 13}

    sc := bufio.NewScanner(f)
    sum := 0
    for sc.Scan() {
        s := sc.Text()
        fmt.Printf("line: %s\n",s)
        id, err := getValidGameID(s, b)
        fmt.Printf("\tvalid: %v\n",id == 0)
        if err == nil {
            sum += id
        }
    }
    fmt.Printf("Final sum: %d\n", sum)
}
