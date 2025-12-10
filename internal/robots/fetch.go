package robots

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Fetcher はrobots.txtを取得・キャッシュする
type Fetcher struct {
	client    *http.Client
	cache     map[string]*RobotsTxt
	userAgent string
}

// NewFetcher は新しいFetcherを作成する
func NewFetcher(userAgent string) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache:     make(map[string]*RobotsTxt),
		userAgent: userAgent,
	}
}

// Fetch は指定されたホストのrobots.txtを取得する
func (f *Fetcher) Fetch(host string, scheme string) (*RobotsTxt, error) {
	// キャッシュを確認
	if cached, ok := f.cache[host]; ok {
		return cached, nil
	}

	// robots.txtのURLを構築
	robotsURL := fmt.Sprintf("%s://%s/robots.txt", scheme, host)

	req, err := http.NewRequest("GET", robotsURL, nil)
	if err != nil {
		// エラー時は空のRobotsTxtを返す（すべて許可）
		empty := &RobotsTxt{}
		f.cache[host] = empty
		return empty, nil
	}

	req.Header.Set("User-Agent", f.userAgent)

	resp, err := f.client.Do(req)
	if err != nil {
		// ネットワークエラー時は空のRobotsTxtを返す
		empty := &RobotsTxt{}
		f.cache[host] = empty
		return empty, nil
	}
	defer resp.Body.Close()

	// 404や他のエラーの場合も空のRobotsTxtを返す
	if resp.StatusCode != http.StatusOK {
		empty := &RobotsTxt{}
		f.cache[host] = empty
		return empty, nil
	}

	// ボディを読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		empty := &RobotsTxt{}
		f.cache[host] = empty
		return empty, nil
	}

	// パースしてキャッシュ
	robotsTxt := Parse(string(body), f.userAgent)
	f.cache[host] = robotsTxt

	return robotsTxt, nil
}
