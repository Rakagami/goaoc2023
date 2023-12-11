package main

import (
	"errors"
	"fmt"

	"github.com/Rakagami/goaoc2023/utils"
)

type Race struct {
    time int
    distance int
}

func magicTriangle(i int, j int) int {
    return i * (j-i)
}

// Good enough for this, lol...
// This can be made more efficient with binary search
func calculateBetterTimes(race Race) int {
    counter := 0
    for i:=0; i<race.time; i++ {
        if magicTriangle(i, race.time) > race.distance {
            counter++
        }
    }
    return counter
}

func main() {
    // ========= Parsing =========
    parseState := 0
    races := []Race{}
    times := []int{}
    distances := []int{}
    parseLines := func(s string) error {
        switch parseState {
        case 0:
            arr, err := utils.ParseInts(s)
            times = arr
            parseState++
            return err
        case 1:
            arr, err := utils.ParseInts(s)
            distances = arr
            parseState++
            return err
        default:
            return errors.New("Unexpected state")
        }
    }

    err := utils.IterateLines("input2.txt", parseLines)

    if err != nil {
        panic(err)
    }

    for i := range times {
        races = append(races, Race{
            time: times[i],
            distance: distances[i],
        })
    }

    fmt.Printf("Races: %v\n", races)

    // ========= Logic =========

    product := 1

    for _, race := range races {
        product *= calculateBetterTimes(race)
    }

    fmt.Printf("Better Time Product: %v\n", product)
}
