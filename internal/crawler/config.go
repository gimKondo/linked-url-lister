package crawler

import (
	"net/url"
	"time"
)

// CrawlConfig はクロール実行時の設定パラメータを保持する
type CrawlConfig struct {
	// StartURL はクロール開始URL
	StartURL *url.URL
	// MaxDepth は最大クロール深度（0=無制限）
	MaxDepth int
	// MinTextLength は最小テキスト文字数
	MinTextLength int
	// Delay はリクエスト間隔
	Delay time.Duration
	// RespectRobotsTxt はrobots.txt遵守フラグ
	RespectRobotsTxt bool
	// OutputJSON はJSON出力フラグ
	OutputJSON bool
	// MaxPages は最大ページ数制限（0=無制限）
	MaxPages int
	// UserAgent はUser-Agentヘッダー
	UserAgent string
	// Verbose は詳細進捗表示フラグ
	Verbose bool
}

// DefaultConfig はデフォルト設定を返す
func DefaultConfig() *CrawlConfig {
	return &CrawlConfig{
		MaxDepth:         0,
		MinTextLength:    500,
		Delay:            100 * time.Millisecond,
		RespectRobotsTxt: true,
		OutputJSON:       false,
		MaxPages:         0,
		UserAgent:        "linked-url-lister/1.0",
		Verbose:          false,
	}
}
