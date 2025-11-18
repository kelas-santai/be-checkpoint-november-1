package tools

import "strings"

func RemoveSpaci(text string) string {

	removeText := strings.ReplaceAll(text, " ", "")

	return removeText
}
