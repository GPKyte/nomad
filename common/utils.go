package common

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strings"
)

func factorial(n int64) (product int64) {
	product = 1
	for i := int64(1); i <= n; i++ {
		product = product * i
	}
	return product
}

func permute(original ...interface{}) [][]interface{} {
	var perms = make([][]interface{}, len(original) /*???*/)

	if len(original) <= 1 {
		perms = append(perms, original)
	}
	/* Iterate through a list of choices
	 * and permute remainder to get all permutations
	 * Hard to explain, but fairly simple in concept */
	for i := 0; i < len(original); i++ {
		remainder := append(original[:i], original[i+1:])
		first := original[i : i+1] /* Get slice of size one */

		for _, ordering := range permute(remainder) {
			p := append(first, ordering)
			perms = append(perms, p)
		}
	}
	return perms
}

/* Sends back two slices of int indices which describe the combos when iterated through */
/* Guaranteed that pairs do not contain same index twice, e.g. (1,2) isValid but (8,8) isNotValid */
/* Order does matter, e.g. (1,3) and (3,1) are distinct */
func pairs(sizeOfCollection int) (A []int, B []int) {
	if sizeOfCollection <= 1 {
		return nil, nil
	}
	var anySmallNumberOfPairs = 12

	// get some n pairs
	for count := 0; count < anySmallNumberOfPairs; {
		a, b := rand.Intn(sizeOfCollection), rand.Intn(sizeOfCollection)
		if a != b {
			A = append(A, a)
			B = append(B, b)
			count++
		}
	}
	/* Alt to get all pairs (in order)
	for left := 0; left < sizeOfCollection; left++ {
		for right := 0; right < sizeOfCollection; right++ {
			if a != b {
				A = append(A, a)
				B = append(B, b)
			}
		}
	}*/
	if len(A) != len(B) {
		// Could note the size of collection and other details like len(A) or B
		panic("Size mismatch in pairs()")
	}

	return A, B
}

type timeAndPlaceList []TimeAndPlace

func (slice timeAndPlaceList) Len() int {
	return len(slice)
}
func (slice timeAndPlaceList) Less(i, j int) bool {
	return slice[i].T.Before(slice[j].T)
}
func (slice timeAndPlaceList) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

/* Do we concern with the potential overhead of passing by value structs? */
/* Because we know not enough and the compiler is very good, let's say no for now */
func (slice timeAndPlaceList) sort(leastToGreatest bool) []TimeAndPlace {
	var total = len(slice)
	var S = make(timeAndPlaceList, 0, total)

	/* Copy slice, intent to refactor this util or delete entirely so come back later */
	for w := range slice {
		S[w] = slice[w]
	}

	sort.Sort(S)
	return S
}

func writeFile(filepath string, data []byte) error {
	var readAndWriteMode os.FileMode = 666 // No point in making it executable too
	return ioutil.WriteFile(filepath, data, readAndWriteMode)
}

// Allowing added delim flexibility to test features of regexp, should not impact default newline char
func readFileByLine(filepath string, delim ...string) []string {
	if delim == nil {
		delim = []string{"\n"}
	}

	newline := regexp.MustCompile(strings.Join(delim, "|"))
	content, err := ioutil.ReadFile(filepath)

	if content == nil || err != nil {
		log(err.Error())
		return nil
	}

	fully := -1
	return newline.Split(string(content), fully)
}

func log(msg string) {
	fmt.Println(msg)
}
