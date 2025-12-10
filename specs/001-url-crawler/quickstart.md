# クイックスタート: linked-url-lister

## インストール

### ソースからビルド

```bash
# リポジトリをクローン
git clone https://github.com/gimKondo/linked-url-lister.git
cd linked-url-lister

# ビルド
go build -o linked-url-lister ./cmd/linked-url-lister

# パスの通った場所に移動（オプション）
sudo mv linked-url-lister /usr/local/bin/
```

### go install

```bash
go install github.com/gimKondo/linked-url-lister/cmd/linked-url-lister@latest
```

## 基本的な使い方

### 1. シンプルなクロール

```bash
linked-url-lister https://handbook.gitlab.com/handbook/
```

出力例:
```
https://handbook.gitlab.com/handbook/
https://handbook.gitlab.com/handbook/about/
https://handbook.gitlab.com/handbook/company/
https://handbook.gitlab.com/handbook/engineering/
...
```

### 2. NotebookLMへのインポート

```bash
# クリップボードにコピー（macOS）
linked-url-lister https://docs.example.com/ | pbcopy

# クリップボードにコピー（Linux with xclip）
linked-url-lister https://docs.example.com/ | xclip -selection clipboard
```

その後、NotebookLMの「ソースを追加」→「ウェブサイト」に貼り付け。

### 3. 深度を制限してクロール

大規模サイトでは深度を制限すると効率的:

```bash
# 2階層までに制限
linked-url-lister --max-depth=2 https://docs.example.com/
```

### 4. JSON出力

プログラムで処理する場合:

```bash
linked-url-lister --json https://docs.example.com/ > urls.json
```

```json
{
  "urls": ["..."],
  "stats": {
    "total_visited": 100,
    "total_filtered": 30,
    "duration_seconds": 45.2
  }
}
```

## よくある使用パターン

### ドキュメントサイトのURL収集

```bash
# 技術ドキュメント
linked-url-lister --min-text=1000 https://docs.example.com/

# ブログ記事
linked-url-lister --max-depth=1 https://blog.example.com/
```

### レート制限を考慮したクロール

```bash
# サーバーに優しいクロール（500ms間隔）
linked-url-lister --delay=500ms https://example.com/

# ページ数を制限
linked-url-lister --max-pages=50 https://example.com/
```

### 進捗を確認しながらクロール

```bash
# 詳細モード
linked-url-lister --verbose https://example.com/ 2>&1 | tee crawl.log
```

## トラブルシューティング

### 「開始URLに接続できません」

- URLが正しいか確認（http:// または https:// を含む）
- ネットワーク接続を確認
- プロキシ設定が必要な場合は `HTTP_PROXY` 環境変数を設定

### 「URLが見つかりません」

- 開始URLが正しいパスを指しているか確認
- サイトがJavaScriptで動的にコンテンツを生成していないか確認
- `--min-text=0` で最小テキスト長を無効化してみる

### クロールが遅い

- `--delay` を小さくする（サーバーに負荷をかけすぎないよう注意）
- `--max-depth` や `--max-pages` で範囲を制限

## 次のステップ

- [CLIインターフェース定義](./contracts/cli-interface.md) で全オプションを確認
- [データモデル](./data-model.md) で内部構造を理解
