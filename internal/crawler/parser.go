package crawler

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// ParseHTML はHTMLからリンクとテキストを抽出する
func ParseHTML(body io.Reader) (links []string, textLength int, err error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, 0, err
	}

	links = []string{}
	var textBuilder strings.Builder

	var extractContent func(*html.Node)
	extractContent = func(n *html.Node) {
		// script, style, noscript タグは除外
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "noscript", "iframe":
				return
			}
		}

		// aタグからhrefを抽出
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := strings.TrimSpace(attr.Val)
					if href != "" && !strings.HasPrefix(href, "#") && !strings.HasPrefix(href, "javascript:") && !strings.HasPrefix(href, "mailto:") {
						links = append(links, href)
					}
					break
				}
			}
		}

		// テキストノードからテキストを抽出
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				if textBuilder.Len() > 0 {
					textBuilder.WriteString(" ")
				}
				textBuilder.WriteString(text)
			}
		}

		// 子ノードを再帰的に処理
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractContent(c)
		}
	}

	extractContent(doc)

	// 連続する空白を正規化
	text := normalizeWhitespace(textBuilder.String())
	textLength = len([]rune(text))

	return links, textLength, nil
}

// normalizeWhitespace は連続する空白を1つに正規化する
func normalizeWhitespace(s string) string {
	var result strings.Builder
	prevSpace := false

	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if !prevSpace {
				result.WriteRune(' ')
				prevSpace = true
			}
		} else {
			result.WriteRune(r)
			prevSpace = false
		}
	}

	return strings.TrimSpace(result.String())
}
