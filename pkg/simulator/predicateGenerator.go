package simulator

import (
	"math/rand"
)

func AdvPredicateGenerator() (string, string, string) {
	subjectList := []string{"apple", "tesla", "microsoft", "amazon", "nvidia"}
	valueList := []string{"50", "80", "85", "90", "99"}

	subject := subjectList[rand.Intn(len(subjectList))]
	operator := ">"
	value := valueList[rand.Intn(len(valueList))]

	return subject, operator, value
}

func SubPredicateGenerator() (string, string, string) {
	subjectList := []string{"apple", "tesla", "microsoft", "amazon", "nvidia"}
	valueList := []string{"100", "130", "150", "170", "190"}

	subject := subjectList[rand.Intn(len(subjectList))]
	operator := ">"
	value := valueList[rand.Intn(len(valueList))]

	return subject, operator, value
}

func PubPredicateGenerator() (string, string, string) {
	subjectList := []string{"apple", "tesla", "microsoft", "amazon", "nvidia"}
	valueList := []string{"100", "123", "147", "166", "182"}

	subject := subjectList[rand.Intn(len(subjectList))]
	operator := "="
	value := valueList[rand.Intn(len(valueList))]

	return subject, operator, value
}
