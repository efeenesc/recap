package db

import "strings"

func generateNumOfQuestionMarks(number int) string {
	questionMarks := []string{}

	for idx := range number {
		_ = idx
		questionMarks = append(questionMarks, "?")
	}

	return strings.Join(questionMarks, ",")
}
