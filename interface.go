package astscan

import (
	"go/ast"
	"regexp"
	"unicode"
)

var (
	reChinesePunctuation = regexp.MustCompile("/·|，|。|《|》|‘|’|”|“|；|：|【|】|？|（|）|、/")
)

// Callback is the scanning callbacker
type Callback func(dir, fileline, astName string, node ast.Node)

// Container is the scanning checker
type Container func(s string) bool

// ContainChinese check whether s contains Chinese character
func ContainChinese(s string) bool {
	chars := unicode.Scripts["Han"]
	for _, r := range s {
		if unicode.Is(chars, r) || reChinesePunctuation.Match([]byte(string(r))) {
			return true
		}
	}
	return false
}
