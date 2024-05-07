package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func NextDirName(basePath string) (string, error) {
	// basePath 내의 모든 항목을 읽음
	items, err := os.ReadDir(basePath)
	if err != nil {
		return "", err
	}

	maxNumber := 0
	for _, item := range items {
		if item.IsDir() {
			// 디렉토리 이름에서 숫자 파싱
			num, err := strconv.Atoi(item.Name())
			if err == nil && num > maxNumber {
				maxNumber = num
			}
		}
	}

	// 가장 높은 숫자에 1을 더함
	newDirName := fmt.Sprintf("%d", maxNumber+1)
	newDirPath := filepath.Join(basePath, newDirName)

	// 새 디렉토리 생성
	err = os.Mkdir(newDirPath, 0755)
	if err != nil {
		return "", err
	}

	return newDirPath, nil
}
