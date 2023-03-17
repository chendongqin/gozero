package words

import (
	"os"
	"path"
	"runtime"
	"strings"
)

type WordsService struct {
	ForbiddenMap map[string]bool //敏感词
	BannedMap    map[string]bool //禁用词
}

func NewWordsService() *WordsService {
	_, filename, _, _ := runtime.Caller(0)
	filePath := path.Dir(filename)
	forbiddenFile := filePath + "/file/forbidden.txt"
	content, err := os.ReadFile(forbiddenFile)
	forbiddenMap := map[string]bool{}
	bannedMap := map[string]bool{}
	if err == nil {
		contentArr := strings.Split(string(content), "\n")
		for _, v := range contentArr {
			forbiddenMap[strings.ToLower(v)] = true
		}
	}
	bannedFile := filePath + "/file/banned.txt"
	content, err = os.ReadFile(bannedFile)
	if err == nil {
		contentArr := strings.Split(string(content), "\n")
		for _, v := range contentArr {
			bannedMap[strings.ToLower(v)] = true
		}
	}
	return &WordsService{
		ForbiddenMap: forbiddenMap,
		BannedMap:    bannedMap,
	}
}

func (receiver *WordsService) CheckAll(content string) string {
	content = strings.ToLower(content)
	if len(receiver.BannedMap) == 0 && len(receiver.ForbiddenMap) == 0 {
		return ""
	}
	wordLen := len(content)
	findWords := make([]string, 0)
	for i := 0; i < wordLen; i++ {
		for j := i + 1; j <= wordLen; j++ {
			subStr := content[i:j]
			if _, found := receiver.BannedMap[subStr]; found {
				findWords = append(findWords, subStr)
				continue
			}
			if _, found := receiver.ForbiddenMap[subStr]; found {
				findWords = append(findWords, subStr)
			}
		}
	}
	return strings.Join(findWords, ",")
}

func (receiver *WordsService) CheckForbidden(content string) string {
	content = strings.ToLower(content)
	if len(receiver.ForbiddenMap) == 0 {
		return ""
	}
	wordLen := len(content)
	findWords := make([]string, 0)
	for i := 0; i < wordLen; i++ {
		for j := i + 1; j <= wordLen; j++ {
			subStr := content[i:j]
			if _, found := receiver.ForbiddenMap[subStr]; found {
				findWords = append(findWords, subStr)
			}
		}
	}
	return strings.Join(findWords, ",")
}

func (receiver *WordsService) CheckBanned(content string) string {
	if len(receiver.BannedMap) == 0 {
		return ""
	}
	wordLen := len(content)
	findWords := make([]string, 0)
	for i := 0; i < wordLen; i++ {
		for j := i + 1; j <= wordLen; j++ {
			subStr := content[i:j]
			if _, found := receiver.BannedMap[subStr]; found {
				findWords = append(findWords, subStr)
			}
		}
	}
	return strings.Join(findWords, ",")
}
