package filter

import (
	"net/url"
	"strings"
)

// IsUnderPath は対象URLが基準URLの配下にあるかを判定する
func IsUnderPath(baseURL, targetURL *url.URL) bool {
	// ホストが一致するか確認
	if baseURL.Host != targetURL.Host {
		return false
	}

	// スキームが一致するか確認（http/https混在は許可しない）
	if baseURL.Scheme != targetURL.Scheme {
		return false
	}

	// パスプレフィックスが一致するか確認
	basePath := baseURL.Path
	targetPath := targetURL.Path

	// トレイリングスラッシュを正規化して比較
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	// ベースパスと完全一致、またはベースパスの配下にあるか
	return targetPath == baseURL.Path || strings.HasPrefix(targetPath, basePath) || targetPath+"/" == basePath
}
