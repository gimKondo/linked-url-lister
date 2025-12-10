package output

import (
	"fmt"
	"os"
)

// Progress は進捗出力を管理する
type Progress struct {
	verbose bool
	count   int
}

// NewProgress は新しいProgressを作成する
func NewProgress(verbose bool) *Progress {
	return &Progress{
		verbose: verbose,
		count:   0,
	}
}

// Found は発見したURLを出力する
func (p *Progress) Found(url string, linkCount int, textLength int) {
	p.count++
	if p.verbose {
		fmt.Fprintf(os.Stderr, "[%d] 発見: %s\n", p.count, url)
		fmt.Fprintf(os.Stderr, "  → %dリンク発見、テキスト長: %d文字\n", linkCount, textLength)
	} else {
		fmt.Fprintf(os.Stderr, "[%d] 発見: %s\n", p.count, url)
	}
}

// Skip はスキップしたURLを出力する
func (p *Progress) Skip(url string, reason string) {
	if p.verbose {
		fmt.Fprintf(os.Stderr, "[スキップ] %s: %s\n", reason, url)
	}
}

// Error はエラーを出力する
func (p *Progress) Error(url string, err string) {
	fmt.Fprintf(os.Stderr, "エラー: %s - %s\n", url, err)
}

// Summary は完了サマリーを出力する
func (p *Progress) Summary(totalFound int, totalVisited int, totalFiltered int, durationSec float64) {
	fmt.Fprintf(os.Stderr, "\n完了: %d件のURLを発見 (訪問: %d, フィルタ: %d, 所要時間: %.1f秒)\n",
		totalFound, totalVisited, totalFiltered, durationSec)
}
