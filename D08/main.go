package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/Rakagami/goaoc2023/utils"
)

type Element string
type ElementSet map[Element]bool

type InstructionElementHash [16]byte
// Returns a canonical Md5 hash of element set and move
func (es Element) Hash(instruction string) InstructionElementHash {
    hash := md5.Sum([]byte(string(es) + instruction))
    return hash
}

func (es ElementSet) Add(e Element) ElementSet {
    es[e] = true
    return es
}

func (es ElementSet) Remove(e Element) ElementSet {
    delete(es, e)
    return es
}

func (es ElementSet) Union(other ElementSet) ElementSet {
    for e, _ := range other {
        es[e] = true
    }
    return es
}

// Checkes whether element set is in final state
func (es ElementSet) IsFinal(finalSet ElementSet) bool {
    for e := range finalSet {
        _, ok := es[e]
        if !ok {
            return false
        }
    }
    return true
}

type Rule struct {
    src Element
    left Element
    right Element
}

type RuleSet struct {
    ruleMap map[Element]Rule  // Map of rules for efficient finding
}

func NewRuleSet() RuleSet {
    return RuleSet{
        ruleMap: map[Element]Rule{},
    }
}

func (rs RuleSet) AddRule(r Rule) {
    rs.ruleMap[r.src] = r
}

func (rs RuleSet) RemoveRule(r Rule) {
    delete(rs.ruleMap, r.src)
}

func (rs RuleSet) Move(left bool, e Element) Element {
    if left {
        return rs.ruleMap[e].left
    } else {
        return rs.ruleMap[e].right
    }
}

func (rs RuleSet) MoveInstruction(instructionString string, e Element) Element {
    for _, c := range instructionString {
        switch(c) {
        case 'L':
            e = rs.Move(true, e)
            break
        case 'R':
            e = rs.Move(false, e)
            break
        default:
            panic(errors.New("Unexpected state reached"))
        }
    }
    return e
}

func (rs RuleSet) MoveSet(left bool, es ElementSet) ElementSet {
    retEs := ElementSet{}
    for e := range es {
        retEs = retEs.Add(rs.Move(left, e))
    }
    return retEs
}

func (rs RuleSet) MoveSetInstruction(instructionString string, es ElementSet, lookupTable map[InstructionElementHash]Element) (ElementSet, map[InstructionElementHash]Element) {
    returnSet := ElementSet{}
    for e := range es {
        lookupUpElement := rs.MoveInstruction(instructionString, e)
        returnSet = returnSet.Add(lookupUpElement)
        lookupTable[e.Hash(instructionString)] = lookupUpElement
    }
    return returnSet, lookupTable
}

// Returns starting and ending ElementSets
func (rs RuleSet) GetTerminalElements() (ElementSet, ElementSet) {
    starting := ElementSet{}
    ending := ElementSet{}

    for el := range rs.ruleMap {
        c := el[2]
        if c == 'A' {
            starting = starting.Add(el)
        } else if c == 'Z' {
            ending = ending.Add(el)
        }
    }

    return starting, ending
}

// Parses the string into a rule
func parseRule(s string) Rule {
    re := regexp.MustCompile("[A-Z0-9]{3}")
    
    matchArr := re.FindAll([]byte(s), 3)

    return Rule{
        src: Element(matchArr[0]),
        left: Element(matchArr[1]),
        right: Element(matchArr[2]),
    }
}

func main() {
    // ============ Parsing ============

    ruleSet := NewRuleSet()
    parseState := 0
    baseInstruction := ""
    parseLines := func(s string) error {
        switch(parseState) {
        case 0:
            if s != "" {
                baseInstruction = s
            } else {
                parseState++
                return nil
            }
        case 1:
            ruleSet.AddRule(parseRule(s))
        }
        return nil
    }

    utils.IterateLines("input_test2.txt", parseLines)
 
    // ============ Logic ============

    //fmt.Printf("RuleSet %v\n", ruleSet)
    fmt.Printf("baseInstruction %v\n", baseInstruction)

    // Current Element
    cur, finalSet := ruleSet.GetTerminalElements()

    fmt.Printf("StartingSet: %v\n", cur)
    fmt.Printf("FinalSet: %v\n", finalSet)
    iterationCnt := 0
    lookupTable := map[InstructionElementHash]Element{}
    for !cur.IsFinal(finalSet) {
        // Check which elements can be found in the lookup table
        lookedUpElements := ElementSet{}
        restElements := ElementSet{}
        for c := range cur {
            hash := c.Hash(baseInstruction)
            _, ok := lookupTable[hash]
            if ok {
                fmt.Printf("\tUsed Lookup table\n")
                lookedUpElements = lookedUpElements.Add(c)
            } else {
                restElements = restElements.Add(c)
            }
        }

        // If there are some unlooked up elements, iterate through
        if len(restElements) > 0 {
            var calculatedElements ElementSet
            calculatedElements, lookupTable = ruleSet.MoveSetInstruction(baseInstruction, restElements, lookupTable)
            lookedUpElements = lookedUpElements.Union(calculatedElements)
        }

        cur = lookedUpElements
        iterationCnt++

        fmt.Printf("cur: %v\n", cur)
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Printf("Counter: %v\n", iterationCnt * len(baseInstruction))
}
