package main

import (
	"fmt"
	"regexp"

	"github.com/Rakagami/goaoc2023/utils"
)

type Element string
type ElementSet map[Element]bool

func (es ElementSet) Add(e Element) ElementSet {
    es[e] = true
    return es
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

func (rs RuleSet) Move(left bool, e Element) Element {
    if left {
        return rs.ruleMap[e].left
    } else {
        return rs.ruleMap[e].right
    }
}

func (rs RuleSet) CountTillTerminal(instructionString string, e Element) int {
    cnt := 0
    for e[2] != 'Z' {
        c := instructionString[cnt % len(instructionString)]
        e = rs.Move(c == 'L', e)
        cnt++
    }
    return cnt
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

    utils.IterateLines("input2.txt", parseLines)
 
    // ============ Logic ============

    fmt.Printf("baseInstruction %v\n", baseInstruction)

    // Current Element
    cur, finalSet := ruleSet.GetTerminalElements()

    fmt.Printf("StartingSet: %v\n", cur)
    fmt.Printf("FinalSet: %v\n", finalSet)

    terminalSteps := []int{}
    for e := range cur {
        terminalSteps = append(terminalSteps, ruleSet.CountTillTerminal(baseInstruction, e))
    }

    fmt.Printf("Smallest Common Multiple: %v\n", utils.ScmArr(terminalSteps))
}
