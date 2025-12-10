package crawler

import "sync"

// VisitedURLs は訪問済みURLを追跡する
type VisitedURLs struct {
	mu      sync.RWMutex
	visited map[string]bool
}

// NewVisitedURLs は新しいVisitedURLsを作成する
func NewVisitedURLs() *VisitedURLs {
	return &VisitedURLs{
		visited: make(map[string]bool),
	}
}

// Add はURLを訪問済みとしてマークする
// すでに訪問済みの場合はfalseを返す
func (v *VisitedURLs) Add(url string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.visited[url] {
		return false
	}

	v.visited[url] = true
	return true
}

// Has はURLが訪問済みかを確認する
func (v *VisitedURLs) Has(url string) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.visited[url]
}

// Count は訪問済みURLの数を返す
func (v *VisitedURLs) Count() int {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return len(v.visited)
}
