package crawler

import "net/url"

// CrawlURL はクロールキューに追加されるURL情報
type CrawlURL struct {
	// URL は正規化された絶対URL
	URL *url.URL
	// Depth は開始URLからの深度（0起点）
	Depth int
	// SourceURL はこのURLを発見した元ページ
	SourceURL string
}
