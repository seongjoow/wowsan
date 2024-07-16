package simulator

import (
	// uuid "github.com/satori/go.uuid"
	"fmt"
	"math/rand"
	"time"
)

func AdvPredicateGenerator(subjectList []string) (string, string, string, []string) {
	// valueList := []string{"50", "60", "70", "80", "90"}
	// valueList := []string{string(uuid.NewV4().String())}

	selectedSubjectList := []string{}

	subject := subjectList[rand.Intn(len(subjectList))]
	selectedSubjectList = append(selectedSubjectList, subject)
	operator := ">"
	// value := valueList[rand.Intn(len(valueList))]
	rand.Seed(time.Now().UnixNano()) // 현재 시간을 기반으로 난수 생성기에 시드값 제공
	// value := rand.Intn(50) + 1       // 1부터 100까지의 랜덤 숫자 생성
	value := rand.Intn(10000000) + 1 // 1부터 100까지의 랜덤 숫자 생성
	valueStr := fmt.Sprint(value)

	return subject, operator, valueStr, selectedSubjectList
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
