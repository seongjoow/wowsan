package simulator

import (
	// uuid "github.com/satori/go.uuid"
	"math/rand"
	"strconv"
)

func AdvPredicateGenerator(subjectList []string) (string, string, string, []string) {
	// valueList := []string{"50", "60", "70", "80", "90"}

	// valueList := []string{string(uuid.NewV4().String())}

	selectedSubjectList := []string{}

	subject := subjectList[rand.Intn(len(subjectList))]
	selectedSubjectList = append(selectedSubjectList, subject)
	operator := ">"
	// value := valueList[rand.Intn(len(valueList))]
	value := strconv.FormatFloat(rand.Float64()*rand.ExpFloat64(), 'f', -1, 64) // Generate a random float, possibly very large

	return subject, operator, value, selectedSubjectList
}

func SubPredicateGenerator(subjectList []string) (string, string, string) {
	// subjectList := []string{"apple", "tesla", "microsoft", "amazon", "nvidia"}
	// valueList := []string{"45", "55", "65", "75", "85", "95"}
	valueList := []string{"45"}

	subject := subjectList[rand.Intn(len(subjectList))]
	operator := ">"
	value := valueList[rand.Intn(len(valueList))]

	return subject, operator, value
}

func PubPredicateGenerator(subjectList []string) (string, string, string) {
	// valueList := []string{"100", "123", "147", "166", "182"}
	// valueList := []string{"100", "110", "120", "130", "140"}
	valueList := []string{"100"}

	subject := subjectList[rand.Intn(len(subjectList))]
	operator := "="
	value := valueList[rand.Intn(len(valueList))]

	return subject, operator, value
}
