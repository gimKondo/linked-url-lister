package crawler

import (
	"fmt"
	"os"
)

// ProgressReporter は進捗出力を管理する
type ProgressReporter struct {
	verbose bool
	count   int
}

// NewProgressReporter は新しいProgressReporterを作成する
func NewProgressReporter(verbose bool) *ProgressReporter {
	return &ProgressReporter{
		verbose: verbose,
		count:   0,
	}
}

// Found は発見したURLを出力する
func (p *ProgressReporter) Found(url string, linkCount int, textLength int) {
	p.count++
	if p.verbose {
		fmt.Fprintf(os.Stderr, "[%d] 発見: %s\n", p.count, url)
		fmt.Fprintf(os.Stderr, "  → %dリンク発見、テキスト長: %d文字\n", linkCount, textLength)
	} else {
		fmt.Fprintf(os.Stderr, "[%d] 発見: %s\n", p.count, url)
	}
}

// Skip はスキップしたURLを出力する
func (p *ProgressReporter) Skip(url string, reason string) {
	if p.verbose {
		fmt.Fprintf(os.Stderr, "[スキップ] %s: %s\n", reason, url)
	}
}

// Error はエラーを出力する
func (p *ProgressReporter) Error(url string, err string) {
	fmt.Fprintf(os.Stderr, "エラー: %s - %s\n", url, err)
}
