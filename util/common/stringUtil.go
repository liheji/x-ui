package common

import (
	"sort"
	"strconv"
	"strings"
)

func IsSubString(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	return index < len(str_array) && str_array[index] == target
}

func CmpVersion(ver1, ver2 string) int {
	splitWithDot := func(ver string) []int {
		newVer := strings.ReplaceAll(ver, "v", "")
		parts := strings.Split(newVer, ".")
		versions := make([]int, 0)
		for _, part := range parts {
			num, _ := strconv.Atoi(part)
			versions = append(versions, num)
		}
		return versions
	}

	v1 := splitWithDot(ver1)
	v2 := splitWithDot(ver2)

	maxLen := len(v1)
	if maxLen < len(v2) {
		maxLen = len(v2)
	}
	for i := 0; i < maxLen; i++ {
		val1 := 0
		val2 := 0

		if i < len(v1) {
			val1 = v1[i]
		}

		if i < len(v2) {
			val2 = v2[i]
		}

		if val1 > val2 {
			return 1
		} else if val1 < val2 {
			return -1
		}
	}

	return 0
}
