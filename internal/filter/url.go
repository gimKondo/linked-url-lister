package filter

import (
	"net/url"
	"path"
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

	// クエリパラメータを除去（同一ページの重複防止）
	resolved.RawQuery = ""

	// パス正規化: 連続スラッシュとドットセグメントを解決
	resolved.Path = path.Clean(resolved.Path)

	// 空パス処理: 空または"."を"/"に変換
	if resolved.Path == "" || resolved.Path == "." {
		resolved.Path = "/"
	}

	// トレイリングスラッシュ除去: ルートパス"/"以外で末尾の"/"を除去
	if resolved.Path != "/" {
		resolved.Path = strings.TrimSuffix(resolved.Path, "/")
	}

	return resolved, nil
}

// IsValidURL はURLが有効なHTTP/HTTPS URLかを判定する
func IsValidURL(u *url.URL) bool {
	return u.Scheme == "http" || u.Scheme == "https"
}
