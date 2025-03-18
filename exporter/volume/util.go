package volume

import (
	"regexp"
	"sort"
	"strconv"
)

const (
	template string = `\d+$`
)

// SortByNumericSuffix 파일명을 숫자 기준으로 정렬하는 함수
func SortByNumericSuffix(files []string) {
	sort.Slice(files, func(i, j int) bool {
		// 숫자 추출
		numI := extractNumber(files[i])
		numJ := extractNumber(files[j])

		// 숫자 기준으로 정렬
		return numI < numJ
	})
}

// extractNumber 파일명에서 숫자 추출 함수
func extractNumber(filename string) int {
	re := regexp.MustCompile(template) // 파일명 끝의 숫자를 찾는 정규식

	match := re.FindString(filename) // 정규식으로 숫자 찾기
	if match == "" {
		return 0
	}
	num, _ := strconv.Atoi(match) // 문자열을 정수로 변환
	return num
}
