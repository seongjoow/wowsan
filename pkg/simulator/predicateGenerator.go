package simulator

import (
	"math/rand"
)

func AdvPredicateGenerator(subjectList []string) (string, string, string, []string) {
	valueList := []string{"50", "80", "85", "90", "99"}
	selectedSubjectList := []string{}

	subject := subjectList[rand.Intn(len(subjectList))]
	selectedSubjectList = append(selectedSubjectList, subject)
	operator := ">"
	value := valueList[rand.Intn(len(valueList))]

	return subject, operator, value, selectedSubjectList
}

func SubPredicateGenerator() (string, string, string) {
	subjectList := []string{"apple", "tesla", "microsoft", "amazon", "nvidia"}
	// valueList := []string{"100", "130", "150", "170", "190"}
	valueList := []string{"60", "70", "95"}

	subject := subjectList[rand.Intn(len(subjectList))]
	operator := ">"
	value := valueList[rand.Intn(len(valueList))]

	return subject, operator, value
}

func PubPredicateGenerator(subjectList []string) (string, string, string) {
	// valueList := []string{"100", "123", "147", "166", "182"}
	valueList := []string{"100", "110", "120"}

	subject := subjectList[rand.Intn(len(subjectList))]
	operator := "="
	value := valueList[rand.Intn(len(valueList))]

	return subject, operator, value
}
