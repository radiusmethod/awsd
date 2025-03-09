package utils

import (
	"github.com/radiusmethod/promptui/list"
	"os"
	"strings"
)

const (
	NoticeColor  = "\033[0;38m%s\u001B[0m"
	PromptColor  = "\033[1;38m%s\u001B[0m"
	CyanColor    = "\033[0;36m%s\033[0m"
	MagentaColor = "\033[0;35m%s\033[0m"
)

type BellSkipper struct{}

func NewPromptUISearcher(items []string) list.Searcher {
	return func(searchInput string, itemIndex int) bool {
		return strings.Contains(strings.ToLower(items[itemIndex]), strings.ToLower(searchInput))
	}
}

func (bs *BellSkipper) Write(b []byte) (int, error) {
	const charBell = 7 // c.f. readline.CharBell
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

func (bs *BellSkipper) Close() error {
	return os.Stderr.Close()
}
