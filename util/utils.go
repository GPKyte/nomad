package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
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

	/* Iterate through a list of choices
	 * and permute remainder to get all permutations
	 * Hard to explain, but fairly simple in concept */
	for i := 0; i < len(original); i++ {
		remainder := append(original[:i], original[i+1:])
		first := original[i : i+1] /* Get slice of size one */

		for ordering := range permute(remainder) {
			p := append(first, ordering)
			perms = append(perms, p)
		}
	}
	return perms
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
