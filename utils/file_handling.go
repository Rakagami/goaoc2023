package utils

import (
	"bufio"
	"os"
)

// Iterate through the lines that takes a callback
func IterateLines(filePath string, callback func(string) error) error {
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

