package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Rakagami/goaoc2023/utils"
)

type Hand string

// Returns tier of hand. Tiers are:
//  * Tier 6: Five of a kind
//  * Tier 5: Four of a kind
//  * Tier 4: Full house
//  * Tier 3: Three of a kind
//  * Tier 2: Two pair
//  * Tier 1: One pair
//  * Tier 0: High Card
// (The higher, the better)
func (hand Hand) Tier() int {
    dict := map[rune]int{}
    for _, c := range hand {
        if val, ok := dict[c]; ok {
            dict[c] = val + 1
        } else {
            dict[c] = 1
        }
    }

    countValue := func(i int) int {
        count := 0
        for _, value := range dict {
            if i == value {
                count++
            }
        }
        return count
    }

    if countValue(5) > 0 {
        return 6
    } else if countValue(4) > 0 {
        return 5
    } else if countValue(3) > 0 && countValue(2) > 0 {
        return 4
    } else if countValue(3) > 0 {
        return 3
    } else if countValue(2) == 2 {
        return 2
    } else if countValue(2) == 1 {
        return 1
    } else {
        return 0
    }
}

func (hand1 Hand) Less(hand2 Hand) bool {
    tier1, tier2 := hand1.Tier(), hand2.Tier()
    if tier1 < tier2 {
        return true
    } else if tier1 > tier2 {
        return false
    }

    labelValue := func(c byte) int {
        if '0' <= c && c <= '9' {
            return int(c - '0')
        }
        switch c {
        case 'T':
            return 10
        case 'J':
            return 11
        case 'Q':
            return 12
        case 'K':
            return 13
        case 'A':
            return 14
        default:
            panic(errors.New("Unexpected label byte"))
        }
    }

    labelLess := func(c1 byte, c2 byte) bool {
        return labelValue(c1) < labelValue(c2)
    }

    // Handling disambiguation
    for i:=0; i<5; i++ {
        c1, c2 := hand1[i], hand2[i]
        if c1 == c2 {
            continue
        } else if labelLess(c1, c2) {
            return true
        } else {
            return false
        }
    }

    panic(errors.New("Mama told me we never reach this state"))
}

type Bid struct {
    hand Hand
    amount int
}

type Bids []Bid
func (bids Bids) Len() int  {return len(bids)}
func (bids Bids) Swap(i, j int) {bids[i], bids[j] = bids[j], bids[i]}
func (bids Bids) Less(i, j int) bool {return bids[i].hand.Less(bids[j].hand)}
func (bids Bids) Sort() {
    sort.Sort(bids)
}

func main() {
    // ============== Parsing ==============

    bids := Bids{}

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

    err := utils.IterateLines("input.txt", parseLines)

    if err != nil {
        panic(err)
    }

    // ============== Logic ==============

    bids.Sort()
    fmt.Printf("Sorted Bids: %v\n", bids)

    sum := 0
    for i, bid := range bids {
        sum += bid.amount * (i+1)
    }

    fmt.Printf("Sum: %v\n", sum)
}
