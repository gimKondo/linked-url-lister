package crawler

import "time"

// CrawlResult はクロール完了後の結果
type CrawlResult struct {
	// URLs はフィルタを通過したURL一覧
	URLs []string
	// TotalVisited は訪問した総ページ数
	TotalVisited int
	// TotalFiltered はフィルタで除外されたページ数
	TotalFiltered int
	// Errors は発生したエラー一覧
	Errors []CrawlError
	// Duration はクロール所要時間
	Duration time.Duration
}
