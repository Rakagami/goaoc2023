package main

import (
	"errors"
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

func (es ElementSet) Remove(e Element) ElementSet {
    delete(es, e)
    return es
}

func (es ElementSet) Terminal() bool {
    _, ok := es["ZZZ"]
    return ok
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

//func (rs RuleSet) Move(left bool, es ElementSet) ElementSet {
//}

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

    utils.IterateLines("input.txt", parseLines)
 
    // ============ Logic ============

    fmt.Printf("RuleSet %v\n", ruleSet)
    fmt.Printf("baseInstruction %v\n", baseInstruction)

    // Current Element
    cur := Element("AAA")
    cnt := 0
    for cur != "ZZZ" {
        instruction := baseInstruction[cnt % len(baseInstruction)]
        switch(instruction) {
        case 'L':
            cur = ruleSet.Move(true, cur)
            cnt++
            break
        case 'R':
            cur = ruleSet.Move(false, cur)
            cnt++
            break
        default:
            panic(errors.New("Unexpected state reached"))
        }
    }

    fmt.Printf("Counter: %v\n", cnt)
}
