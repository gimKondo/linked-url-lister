# CLIインターフェース定義: linked-url-lister

**バージョン**: 1.0.0
**日付**: 2025-12-10

## コマンド概要

```
linked-url-lister [オプション] <URL>
```

指定したURLを起点に、配下のページURLをクロールしてリストアップする。

## 引数

| 引数 | 必須 | 説明 |
|------|------|------|
| URL | ✅ | クロール開始URL（http:// または https://） |

## オプション

| フラグ | 短縮形 | 型 | デフォルト | 説明 |
|--------|--------|-----|----------|------|
| --max-depth | -d | int | 0 | 最大クロール深度（0=無制限） |
| --min-text | -t | int | 500 | 最小テキスト文字数 |
| --delay | | duration | 100ms | リクエスト間隔 |
| --max-pages | -m | int | 0 | 最大ページ数（0=無制限） |
| --json | -j | bool | false | JSON形式で出力 |
| --ignore-robots | | bool | false | robots.txtを無視 |
| --user-agent | -u | string | linked-url-lister/1.0 | User-Agentヘッダー |
| --verbose | -v | bool | false | 詳細な進捗表示 |
| --help | -h | bool | | ヘルプを表示 |
| --version | -V | bool | | バージョンを表示 |

## 使用例

### 基本的な使用法

```bash
# GitLab Handbookの配下URLをリストアップ
linked-url-lister https://handbook.gitlab.com/handbook/

# 深度2までに制限してクロール
linked-url-lister --max-depth=2 https://docs.example.com/

# 1000文字以上のテキストを持つページのみ
linked-url-lister --min-text=1000 https://blog.example.com/
```

### JSON出力

```bash
# JSON形式で出力
linked-url-lister --json https://docs.example.com/

# ファイルに保存
linked-url-lister --json https://docs.example.com/ > urls.json
```

### レート制限

```bash
# 500msの間隔でリクエスト
linked-url-lister --delay=500ms https://example.com/

# 100ページまでに制限
linked-url-lister --max-pages=100 https://example.com/
```

### パイプライン連携

```bash
# 発見したURLをxargsで処理
linked-url-lister https://docs.example.com/ | xargs -I{} curl -s {}

# NotebookLMにコピー可能な形式で出力（デフォルト）
linked-url-lister https://handbook.gitlab.com/handbook/ | pbcopy
```

## 出力

### 標準出力（stdout）

**テキスト形式（デフォルト）**:
```
https://example.com/page1
https://example.com/page2
https://example.com/docs/intro
```

**JSON形式（--json）**:
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
  "errors": []
}
```

### 標準エラー出力（stderr）

進捗情報とエラーメッセージ:

```
[1] 発見: https://example.com/
[2] 発見: https://example.com/page1
[3] スキップ（テキスト不足）: https://example.com/images
エラー: https://example.com/broken - 接続タイムアウト
```

詳細モード（--verbose）:
```
[1/50] クロール中: https://example.com/
  → 10リンク発見、テキスト長: 2500文字
[2/50] クロール中: https://example.com/page1
  → 5リンク発見、テキスト長: 1800文字
```

## 終了コード

| コード | 説明 |
|--------|------|
| 0 | 正常終了 |
| 1 | 引数エラー（無効なURL、不正なオプション値） |
| 2 | ネットワークエラー（開始URLに接続できない） |
| 3 | 中断（Ctrl+C） |

## 環境変数

| 変数 | 説明 |
|------|------|
| HTTP_PROXY | HTTPプロキシ（標準のGo HTTPクライアント設定） |
| HTTPS_PROXY | HTTPSプロキシ |
| NO_PROXY | プロキシを使用しないホスト |

## 制限事項

- JavaScript実行が必要なSPA（Single Page Application）はサポートしない
- 認証が必要なページは自動的にスキップされる
- 同一ホスト・同一パス配下のページのみがクロール対象
- 1回の実行で複数の開始URLを指定することはできない
