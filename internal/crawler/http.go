package crawler

import (
	"net/http"
	"time"
)

// NewHTTPClient はクロール用のHTTPクライアントを作成する
func NewHTTPClient(userAgent string) *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 最大10回のリダイレクトを許可
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			// リダイレクト先でもUser-Agentを設定
			req.Header.Set("User-Agent", userAgent)
			return nil
		},
	}
}

// FetchPage は指定されたURLからページを取得する
func FetchPage(client *http.Client, urlStr string, userAgent string) (*Page, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return &Page{
			URL:   urlStr,
			Error: err,
		}, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ja,en-US;q=0.9,en;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return &Page{
			URL:   urlStr,
			Error: err,
		}, err
	}
	defer resp.Body.Close()

	page := &Page{
		URL:        urlStr,
		StatusCode: resp.StatusCode,
	}

	// 200以外のステータスコードはエラーとして扱う
	if resp.StatusCode != http.StatusOK {
		return page, nil
	}

	return page, nil
}
