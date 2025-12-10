package crawler

// Page は取得したページの情報
type Page struct {
	// URL はページのURL
	URL string
	// StatusCode はHTTPステータスコード
	StatusCode int
	// TextLength は抽出したテキストの文字数
	TextLength int
	// Links はページ内で発見したリンク
	Links []string
	// Error は取得時のエラー（あれば）
	Error error
}
