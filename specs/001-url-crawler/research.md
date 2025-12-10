# リサーチ: NotebookLMインポート用URLクローラー

**日付**: 2025-12-10
**ブランチ**: `001-url-crawler`

## 調査項目

### 1. GoでのHTMLパース

**決定**: 標準ライブラリ `golang.org/x/net/html` を使用

**理由**:
- Go公式の準標準ライブラリ（golang.org/x配下）
- 外部依存を最小限に抑える憲法原則に準拠
- ストリーミングパーサーでメモリ効率が良い
- HTMLの不正な構造にも耐性がある

**検討した代替案**:
- `goquery`: jQueryライクなAPIで使いやすいが、外部依存が増える
- 正規表現: 信頼性が低く、HTMLパースには不適切

### 2. robots.txt解析

**決定**: 自前で簡易実装

**理由**:
- robots.txtの基本ルール（User-agent, Disallow, Allow）は単純
- 外部ライブラリは過剰な機能を持つことが多い
- 憲法の「シンプルさ優先」原則に従う
- 必要な機能: 特定パスがクロール許可されているかの判定のみ

**検討した代替案**:
- `temoto/robotstxt`: 機能は豊富だが外部依存
- Google's robotstxt: C++実装でGoから使いにくい

### 3. テキストコンテンツ抽出

**決定**: HTMLタグを除去し、可視テキストの文字数をカウント

**理由**:
- `<script>`, `<style>`, `<noscript>` タグ内のテキストは除外
- 連続する空白は1つに正規化
- シンプルな文字数カウントで十分（500文字閾値）

**実装方針**:
```go
// 除外するタグ
excludeTags := []string{"script", "style", "noscript", "iframe"}

// テキスト抽出後、連続空白を正規化してlen()でカウント
```

### 4. URL正規化

**決定**: 標準ライブラリ `net/url` で正規化

**理由**:
- 相対URLを絶対URLに変換
- フラグメント（#以降）を除去
- クエリパラメータは保持（同一ページの異なるビューの可能性）
- トレイリングスラッシュの正規化

**実装方針**:
```go
// base: 現在のページURL
// href: 抽出したリンク
parsed, _ := url.Parse(href)
resolved := base.ResolveReference(parsed)
resolved.Fragment = "" // フラグメント除去
```

### 5. 同一パス配下の判定

**決定**: URLパスのプレフィックスマッチ

**理由**:
- 開始URL `https://example.com/docs/` の場合
- `https://example.com/docs/guide/` → 含める
- `https://example.com/blog/` → 除外
- ホスト名も一致する必要がある

**実装方針**:
```go
func isUnderPath(baseURL, targetURL *url.URL) bool {
    if baseURL.Host != targetURL.Host {
        return false
    }
    return strings.HasPrefix(targetURL.Path, baseURL.Path)
}
```

### 6. 認証検出

**決定**: HTTPステータスコードとリダイレクト先で判定

**理由**:
- 401/403: 明示的な認証エラー
- 302/307でログインページへリダイレクト: 認証必要と判断
- 200以外のレスポンスは基本的にスキップ

**実装方針**:
- `http.Client` の `CheckRedirect` でリダイレクト先を監視
- ログインページのパターン（/login, /signin, /auth）を検出

### 7. 並行処理

**決定**: シングルスレッドで順次処理（Phase 1）

**理由**:
- 100msのリクエスト間隔があるため、並行処理のメリットは限定的
- サーバーへの負荷を考慮
- 複雑性を避ける（憲法原則）
- 将来の拡張として並行処理は検討可能

### 8. 進捗表示

**決定**: stderrに簡潔な進捗メッセージを出力

**理由**:
- stdoutは結果出力専用（パイプライン対応）
- 形式: `[N/M] クロール中: URL` または `[N] 発見: URL`
- 総数が不明な場合は発見数のみ表示

## 技術スタック確定

| 項目 | 選択 | 根拠 |
|------|------|------|
| 言語 | Go 1.21+ | 標準ライブラリの充実、クロスコンパイル容易 |
| HTMLパース | golang.org/x/net/html | 準標準、堅牢 |
| HTTP | net/http | 標準ライブラリ |
| URL処理 | net/url | 標準ライブラリ |
| JSON | encoding/json | 標準ライブラリ |
| robots.txt | 自前実装 | シンプルさ優先 |
| 並行処理 | なし（Phase 1） | YAGNI原則 |

## 未解決事項

なし - すべての技術選択が確定。
