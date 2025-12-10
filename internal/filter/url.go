package filter

import (
	"net/url"
	"strings"
)

// NormalizeURL は相対URLを絶対URLに変換し、フラグメントを除去する
func NormalizeURL(baseURL *url.URL, href string) (*url.URL, error) {
	// hrefをパース
	parsed, err := url.Parse(href)
	if err != nil {
		return nil, err
	}

	// 相対URLを絶対URLに解決
	resolved := baseURL.ResolveReference(parsed)

	// フラグメントを除去
	resolved.Fragment = ""

	// トレイリングスラッシュを正規化（ディレクトリの場合）
	// ただし、ファイル拡張子がある場合はスラッシュを追加しない
	if !strings.Contains(resolved.Path, ".") && !strings.HasSuffix(resolved.Path, "/") && resolved.Path != "" {
		// パスがディレクトリの可能性がある場合、スラッシュは追加しない
		// サーバーの応答に依存するため、そのまま
	}

	return resolved, nil
}

// IsValidURL はURLが有効なHTTP/HTTPS URLかを判定する
func IsValidURL(u *url.URL) bool {
	return u.Scheme == "http" || u.Scheme == "https"
}
