package utils

// Greatest common divisor using the euclidean algorithm
func Gcd(a int, b int) int {
    if a < 0 || b < 0 {
        return -1
    }

    if a < b {
        return Gcd(b, a)
    } else if b == 0 {
        return a
    } else {
        return Gcd(b, a % b)
    }
}

func GcdArr(arr []int) int {
    if len(arr) < 2 {
        return -1
    }

    cur := arr[0]
    for i:=1; i<len(arr); i++ {
        cur = Gcd(cur, arr[i])
    }

    return cur
}

// Smallest common multiple using the euclidean algorithm
func Scm(a int, b int) int {
    return a*b / Gcd(a, b)
}

// Smallest common multiple using the euclidean algorithm
func ScmArr(arr []int) int {
    if len(arr) < 2 {
        return -1
    }

    cur := arr[0]
    for i:=1; i<len(arr); i++ {
        cur = Scm(cur, arr[i])
    }

    return cur
}
