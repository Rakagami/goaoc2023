package utils

import (
	"regexp"
	"strconv"
)

func ParseInts(s string) ([]int, error) {
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
