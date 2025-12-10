package robots

import (
	"bufio"
	"strings"
	"time"
)

// RobotsTxt はrobots.txtの解析結果
type RobotsTxt struct {
	// DisallowPaths は禁止パスのリスト
	DisallowPaths []string
	// AllowPaths は許可パスのリスト
	AllowPaths []string
	// CrawlDelay は指定されたクロール遅延
	CrawlDelay time.Duration
}

// Parse はrobots.txtの内容を解析する
func Parse(content string, userAgent string) *RobotsTxt {
	result := &RobotsTxt{
		DisallowPaths: []string{},
		AllowPaths:    []string{},
		CrawlDelay:    0,
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	currentUserAgent := ""
	isRelevant := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// コメントと空行をスキップ
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// User-agent行の処理
		if strings.HasPrefix(strings.ToLower(line), "user-agent:") {
			currentUserAgent = strings.TrimSpace(strings.TrimPrefix(strings.ToLower(line), "user-agent:"))
			// *（すべて）または指定されたUser-Agentに一致するか確認
			isRelevant = currentUserAgent == "*" || strings.Contains(strings.ToLower(userAgent), currentUserAgent)
			continue
		}

		if !isRelevant {
			continue
		}

		// Disallow行の処理
		if strings.HasPrefix(strings.ToLower(line), "disallow:") {
			path := strings.TrimSpace(strings.TrimPrefix(line, line[:9]))
			if path != "" {
				result.DisallowPaths = append(result.DisallowPaths, path)
			}
			continue
		}

		// Allow行の処理
		if strings.HasPrefix(strings.ToLower(line), "allow:") {
			path := strings.TrimSpace(strings.TrimPrefix(line, line[:6]))
			if path != "" {
				result.AllowPaths = append(result.AllowPaths, path)
			}
			continue
		}

		// Crawl-delay行の処理
		if strings.HasPrefix(strings.ToLower(line), "crawl-delay:") {
			delayStr := strings.TrimSpace(strings.TrimPrefix(line, line[:12]))
			var delaySeconds float64
			if _, err := parseFloat(delayStr, &delaySeconds); err == nil && delaySeconds > 0 {
				result.CrawlDelay = time.Duration(delaySeconds * float64(time.Second))
			}
		}
	}

	return result
}

// parseFloat は文字列を浮動小数点数に変換する（簡易実装）
func parseFloat(s string, result *float64) (bool, error) {
	var val float64
	for i, c := range s {
		if c >= '0' && c <= '9' {
			val = val*10 + float64(c-'0')
		} else if c == '.' {
			// 小数点以下を処理
			decimal := 0.1
			for _, dc := range s[i+1:] {
				if dc >= '0' && dc <= '9' {
					val += float64(dc-'0') * decimal
					decimal /= 10
				} else {
					break
				}
			}
			break
		} else {
			break
		}
	}
	*result = val
	return true, nil
}

// IsAllowed は指定されたパスがクロール許可されているかを判定する
func (r *RobotsTxt) IsAllowed(path string) bool {
	// Allowが明示的に指定されている場合は許可
	for _, allowPath := range r.AllowPaths {
		if strings.HasPrefix(path, allowPath) {
			return true
		}
	}

	// Disallowに一致する場合は禁止
	for _, disallowPath := range r.DisallowPaths {
		if strings.HasPrefix(path, disallowPath) {
			return false
		}
	}

	// デフォルトは許可
	return true
}
