# データモデル: NotebookLMインポート用URLクローラー

**日付**: 2025-12-10
**ブランチ**: `001-url-crawler`

## エンティティ定義

### CrawlConfig（クロール設定）

クロール実行時の設定パラメータを保持する。

| フィールド | 型 | デフォルト | 説明 |
|-----------|-----|----------|------|
| StartURL | string | (必須) | クロール開始URL |
| MaxDepth | int | 0 (無制限) | 最大クロール深度 |
| MinTextLength | int | 500 | 最小テキスト文字数 |
| Delay | time.Duration | 100ms | リクエスト間隔 |
| RespectRobotsTxt | bool | true | robots.txt遵守フラグ |
| OutputJSON | bool | false | JSON出力フラグ |
| MaxPages | int | 0 (無制限) | 最大ページ数制限 |
| UserAgent | string | "linked-url-lister/1.0" | User-Agentヘッダー |

**バリデーションルール**:
- StartURL: 有効なHTTP/HTTPS URLであること
- MaxDepth: 0以上の整数
- MinTextLength: 0以上の整数
- Delay: 0以上

### CrawlURL（クロール対象URL）

クロールキューに追加されるURL情報。

| フィールド | 型 | 説明 |
|-----------|-----|------|
| URL | *url.URL | 正規化された絶対URL |
| Depth | int | 開始URLからの深度（0起点） |
| SourceURL | string | このURLを発見した元ページ |

### Page（ページ情報）

取得したページの情報。

| フィールド | 型 | 説明 |
|-----------|-----|------|
| URL | string | ページのURL |
| StatusCode | int | HTTPステータスコード |
| TextLength | int | 抽出したテキストの文字数 |
| Links | []string | ページ内で発見したリンク |
| Error | error | 取得時のエラー（あれば） |

### CrawlResult（クロール結果）

クロール完了後の結果。

| フィールド | 型 | 説明 |
|-----------|-----|------|
| URLs | []string | フィルタを通過したURL一覧 |
| TotalVisited | int | 訪問した総ページ数 |
| TotalFiltered | int | フィルタで除外されたページ数 |
| Errors | []CrawlError | 発生したエラー一覧 |
| Duration | time.Duration | クロール所要時間 |

### CrawlError（クロールエラー）

クロール中に発生したエラー情報。

| フィールド | 型 | 説明 |
|-----------|-----|------|
| URL | string | エラーが発生したURL |
| Error | string | エラーメッセージ |
| StatusCode | int | HTTPステータスコード（該当する場合） |

### RobotsTxt（robots.txt情報）

robots.txtの解析結果。

| フィールド | 型 | 説明 |
|-----------|-----|------|
| DisallowPaths | []string | 禁止パスのリスト |
| AllowPaths | []string | 許可パスのリスト |
| CrawlDelay | time.Duration | 指定されたクロール遅延 |

## 状態遷移

### URLの状態

```
[未発見] → [キュー追加] → [処理中] → [完了/エラー]
                ↑
                └── [スキップ] (robots.txt禁止、深度超過、パス外)
```

### クロール全体の状態

```
[初期化] → [クロール中] → [完了]
              ↓
           [中断] (ユーザー停止、致命的エラー)
```

## データフロー

```
入力: StartURL
  ↓
[URLキュー] ← 新規URL追加
  ↓
[フェッチ] → Page
  ↓
[フィルタ]
  ├── パス外 → スキップ
  ├── 認証必要 → スキップ
  ├── テキスト不足 → スキップ
  └── 合格 → CrawlResult.URLsに追加
  ↓
[リンク抽出] → [URLキュー]
  ↓
(キューが空になるまで繰り返し)
  ↓
出力: CrawlResult
```

## JSON出力形式

```json
{
  "urls": [
    "https://example.com/page1",
    "https://example.com/page2"
  ],
  "stats": {
    "total_visited": 50,
    "total_filtered": 20,
    "duration_seconds": 120.5
  },
  "errors": [
    {
      "url": "https://example.com/broken",
      "error": "connection timeout",
      "status_code": 0
    }
  ]
}
```

## テキスト出力形式

```
https://example.com/page1
https://example.com/page2
https://example.com/page3
```

（1行に1URL、末尾改行あり）
