package crawler

// CrawlError はクロール中に発生したエラー情報
type CrawlError struct {
	// URL はエラーが発生したURL
	URL string
	// Error はエラーメッセージ
	Error string
	// StatusCode はHTTPステータスコード（該当する場合）
	StatusCode int
}
