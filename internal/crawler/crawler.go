package crawler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gimKondo/linked-url-lister/internal/filter"
	"github.com/gimKondo/linked-url-lister/internal/robots"
)

// Crawler はWebクローラー
type Crawler struct {
	config        *CrawlConfig
	client        *http.Client
	visited       *VisitedURLs
	queue         []*CrawlURL
	result        *CrawlResult
	robotsFetcher *robots.Fetcher
	progress      *ProgressReporter
}

// New は新しいクローラーを作成する
func New(config *CrawlConfig) *Crawler {
	return &Crawler{
		config:        config,
		client:        NewHTTPClient(config.UserAgent),
		visited:       NewVisitedURLs(),
		queue:         []*CrawlURL{},
		result:        &CrawlResult{URLs: []string{}, Errors: []CrawlError{}},
		robotsFetcher: robots.NewFetcher(config.UserAgent),
		progress:      NewProgressReporter(config.Verbose),
	}
}

// Run はクロールを実行する
func (c *Crawler) Run() (*CrawlResult, error) {
	startTime := time.Now()

	// 開始URLをキューに追加
	c.queue = append(c.queue, &CrawlURL{
		URL:       c.config.StartURL,
		Depth:     0,
		SourceURL: "",
	})

	// キューが空になるまで処理
	for len(c.queue) > 0 {
		// 最大ページ数制限のチェック
		if c.config.MaxPages > 0 && c.result.TotalVisited >= c.config.MaxPages {
			break
		}

		// キューから取り出し（FIFO）
		current := c.queue[0]
		c.queue = c.queue[1:]

		// すでに訪問済みならスキップ
		urlStr := current.URL.String()
		if !c.visited.Add(urlStr) {
			continue
		}

		// 深度制限のチェック
		if c.config.MaxDepth > 0 && current.Depth > c.config.MaxDepth {
			continue
		}

		// パスプレフィックスのチェック
		if !filter.IsUnderPath(c.config.StartURL, current.URL) {
			continue
		}

		// robots.txtのチェック
		if c.config.RespectRobotsTxt && !c.isAllowedByRobots(current.URL) {
			c.progress.Skip(urlStr, "robots.txtで禁止")
			continue
		}

		// ページを取得
		page := c.fetchAndParse(current)
		c.result.TotalVisited++

		// エラーチェック
		if page.Error != nil {
			c.result.Errors = append(c.result.Errors, CrawlError{
				URL:        urlStr,
				Error:      page.Error.Error(),
				StatusCode: page.StatusCode,
			})
			c.progress.Error(urlStr, page.Error.Error())
			continue
		}

		// ステータスコードチェック
		if page.StatusCode != http.StatusOK {
			c.result.TotalFiltered++
			c.progress.Skip(urlStr, fmt.Sprintf("ステータス %d", page.StatusCode))
			continue
		}

		// テキスト長チェック
		if !filter.HasEnoughText(page.TextLength, c.config.MinTextLength) {
			c.result.TotalFiltered++
			c.progress.Skip(urlStr, fmt.Sprintf("テキスト不足 (%d文字)", page.TextLength))
			continue
		}

		// 結果に追加
		c.result.URLs = append(c.result.URLs, urlStr)
		c.progress.Found(urlStr, len(page.Links), page.TextLength)

		// リンクをキューに追加
		c.addLinksToQueue(current, page.Links)

		// リクエスト間隔
		if c.config.Delay > 0 {
			time.Sleep(c.config.Delay)
		}
	}

	c.result.Duration = time.Since(startTime)
	return c.result, nil
}

// fetchAndParse はページを取得してパースする
func (c *Crawler) fetchAndParse(crawlURL *CrawlURL) *Page {
	urlStr := crawlURL.URL.String()

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return &Page{URL: urlStr, Error: err}
	}

	req.Header.Set("User-Agent", c.config.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := c.client.Do(req)
	if err != nil {
		return &Page{URL: urlStr, Error: err}
	}
	defer resp.Body.Close()

	page := &Page{
		URL:        urlStr,
		StatusCode: resp.StatusCode,
	}

	if resp.StatusCode != http.StatusOK {
		return page
	}

	// ボディを読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		page.Error = err
		return page
	}

	// HTMLをパース
	links, textLength, err := ParseHTML(bytes.NewReader(body))
	if err != nil {
		page.Error = err
		return page
	}

	page.Links = links
	page.TextLength = textLength
	return page
}

// addLinksToQueue はリンクをキューに追加する
func (c *Crawler) addLinksToQueue(current *CrawlURL, links []string) {
	for _, href := range links {
		normalized, err := filter.NormalizeURL(current.URL, href)
		if err != nil {
			continue
		}

		if !filter.IsValidURL(normalized) {
			continue
		}

		urlStr := normalized.String()
		if c.visited.Has(urlStr) {
			continue
		}

		if !filter.IsUnderPath(c.config.StartURL, normalized) {
			continue
		}

		c.queue = append(c.queue, &CrawlURL{
			URL:       normalized,
			Depth:     current.Depth + 1,
			SourceURL: current.URL.String(),
		})
	}
}

// isAllowedByRobots はrobots.txtでクロールが許可されているかを確認する
func (c *Crawler) isAllowedByRobots(u *url.URL) bool {
	robotsTxt, _ := c.robotsFetcher.Fetch(u.Host, u.Scheme)
	return robotsTxt.IsAllowed(u.Path)
}

