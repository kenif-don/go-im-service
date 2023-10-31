package util

import "sort"

func IndexOfString(target string, arr []string) int {
	return sort.SearchStrings(arr, target)
}
