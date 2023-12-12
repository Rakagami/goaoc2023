package utils

import (
	"fmt"
	"testing"
)

func TestGcd(t *testing.T) {
    a := 28
    b := 21
    gcd := Gcd(a, b)

    if a % gcd != 0 || b % gcd != 0 {
        t.Fatalf("Created GCD is not correct")
    } else if gcd != 7 {
        t.Fatalf("Created GCD is not correct")
    }

    fmt.Printf("Test\n")
}

func TestGcdArr(t *testing.T) {
    arr1 := []int{3}
    arr2 := []int{12, 88}
    arr3 := []int{21, 14, 56}

    gcd1 := GcdArr(arr1)
    gcd2 := GcdArr(arr2)
    gcd3 := GcdArr(arr3)

    if gcd1 != -1 {
        t.Fatalf("Created GCD is not correct")
    } else if gcd2 != 4 {
        t.Fatalf("Created GCD is not correct")
    } else if gcd3 != 7 {
        t.Fatalf("Created GCD is not correct")
    }

    fmt.Printf("Test\n")
}
