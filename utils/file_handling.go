package utils

import (
	"bufio"
	"os"
)

// Iterate through the lines that takes a callback
func IterateLines(filePath string, callback func(string) error) error {
    f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
    if err != nil {
        return err
    }
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

