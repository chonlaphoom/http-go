package main

import (
	"strings"
)

func stringReplaceAll(target string, profaneList []string) string {
	splitBy := " "
	lower := strings.ToLower(target)
	lower_splitted := strings.Split(lower, splitBy)
	target_splitted := strings.Split(target, splitBy)

	var indexToReplace []int
	for _, profane := range profaneList {
		for index, eachSplitted := range lower_splitted {
			if eachSplitted == profane {
				indexToReplace = append(indexToReplace, index)
			}
		}
	}

	for _, index := range indexToReplace {
		target_splitted[index] = "****"
	}

	return strings.Join(target_splitted, " ")
}
