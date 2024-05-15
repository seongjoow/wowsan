package automation

// import (
// 	"fmt"
// 	"io/fs"
// 	"os"
// 	"path/filepath"
// )

// func CheckTenFiles(dirPath string) bool {
// 	var fileCount int
// 	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err // 파일 탐색 중 오류 발생 시 반환
// 		}
// 		if !d.IsDir() {
// 			fileCount++
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		fmt.Println("Error walking through directory:", err)
// 		return false
// 	}
// 	fmt.Println("Number of files in the directory:", fileCount)
// 	return fileCount == 9 // 파일 수가 10개인지 검사
// }

// func RemoveFiles(dirPath string) {
// 	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err // 파일 탐색 중 오류 발생 시 반환
// 		}
// 		if !d.IsDir() {
// 			err := os.Remove(path)
// 			if err != nil {
// 				fmt.Println("Error removing file:", err)
// 			}
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		fmt.Println("Error walking through directory:", err)
// 	}
// }
