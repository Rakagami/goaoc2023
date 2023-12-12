package main

import (
	"fmt"
	"strings"
    "strconv"

	"github.com/Rakagami/goaoc2023/utils"
)

type Hand string

type Bid struct {
    hand Hand
    amount int
}

func main() {
    // ============== Parsing ==============

    bids := []Bid{}

    parseLines := func(s string) error {
        v := strings.SplitN(s, " ", 2)
        amount, err := strconv.Atoi(v[1])
        if err != nil {
            return err
        }
        bids = append(bids, Bid{
            hand: Hand(v[0]),
            amount: amount,
        })
        return nil
    }

    err := utils.IterateLines("input_test.txt", parseLines)

    if err != nil {
        panic(err)
    }

    // ============== Logic ==============

    fmt.Printf("Bids: %v\n", bids)
}
