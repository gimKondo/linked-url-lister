package filter

// HasEnoughText はテキスト長が最小要件を満たすかを判定する
func HasEnoughText(textLength int, minTextLength int) bool {
	return textLength >= minTextLength
}
