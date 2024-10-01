package db

import "strings"

// Helper function for use in SQL queries where a variable number of parameters must be passed
func generateNumOfQuestionMarks(number int) string {
	questionMarks := []string{}

	for idx := range number {
		_ = idx
		questionMarks = append(questionMarks, "?")
	}

	return strings.Join(questionMarks, ",")
}
