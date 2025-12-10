# linked-url-lister

指定したWebサイトの配下ページURLをリストアップするCLIツールです。主な用途はNotebookLMへのURL一括インポートです。

## 機能

- 指定URLから再帰的にリンクを辿り、配下ページのURLを収集
- パスプレフィックスによるフィルタリング（同一ドメイン・配下パスのみ）
- テキストコンテンツ量によるフィルタリング
- robots.txt遵守
- テキスト形式またはJSON形式での出力

## インストール

### ソースからビルド

```bash
go install github.com/gimKondo/linked-url-lister/cmd/linked-url-lister@latest
```

または

```bash
git clone https://github.com/gimKondo/linked-url-lister.git
cd linked-url-lister
go build -o linked-url-lister ./cmd/linked-url-lister
```

## 使用方法

### 基本的な使用例

```bash
# 基本的な使用
linked-url-lister https://docs.example.com/

# 深度制限付き
linked-url-lister --max-depth=2 https://docs.example.com/

# JSON形式で出力
linked-url-lister --json https://docs.example.com/ > urls.json

# 詳細な進捗表示
linked-url-lister --verbose https://docs.example.com/
```

### NotebookLMへのインポート

```bash
# クリップボードにコピー（macOS）
linked-url-lister https://docs.example.com/ | pbcopy

# ファイルに保存
linked-url-lister https://docs.example.com/ > urls.txt
```

## オプション

| オプション | 短縮形 | デフォルト | 説明 |
|-----------|--------|----------|------|
| --min-text | -t | 500 | 最小テキスト文字数 |
| --max-depth | -d | 0 | 最大クロール深度（0=無制限） |
| --delay | | 100ms | リクエスト間隔 |
| --max-pages | -m | 0 | 最大ページ数（0=無制限） |
| --json | -j | false | JSON形式で出力 |
| --ignore-robots | | false | robots.txtを無視 |
| --user-agent | -u | linked-url-lister/1.0 | User-Agentヘッダー |
| --verbose | -v | false | 詳細な進捗表示 |
| --help | -h | | ヘルプを表示 |
| --version | -V | | バージョンを表示 |

## 出力形式

### テキスト形式（デフォルト）

```
https://docs.example.com/getting-started
https://docs.example.com/api/overview
https://docs.example.com/guides/quickstart
```

### JSON形式

```json
{
  "urls": [
    "https://docs.example.com/getting-started",
    "https://docs.example.com/api/overview"
  ],
  "stats": {
    "total_visited": 50,
    "total_filtered": 20,
    "duration_seconds": 30.5
  },
  "errors": []
}
```

## 終了コード

| コード | 説明 |
|--------|------|
| 0 | 正常終了 |
| 1 | 引数エラー |
| 2 | ネットワークエラー |
| 3 | 中断 |

## ライセンス

MIT License
