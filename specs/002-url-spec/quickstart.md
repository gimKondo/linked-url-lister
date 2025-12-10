# クイックスタート: URL正規化の強化

## 概要

この機能はURLクローラーの内部動作を改善するものであり、CLIインターフェースへの変更はありません。ユーザーは既存のコマンドをそのまま使用できます。

## 変更点

### 改善前の動作

クエリパラメータやトレイリングスラッシュの違いにより、同一ページが複数回出力される可能性がありました。

```bash
# 同一ページが異なるURLで出力される例
$ linked-url-lister https://example.com/docs/
https://example.com/docs/guide
https://example.com/docs/guide/
https://example.com/docs/guide?utm_source=twitter
https://example.com/docs/guide?ref=sidebar
```

### 改善後の動作

URL正規化により、同一ページは1回のみ出力されます。

```bash
# 正規化後は重複なし
$ linked-url-lister https://example.com/docs/
https://example.com/docs/guide
```

## 正規化ルール

| 入力URL | 出力URL | 説明 |
|---------|---------|------|
| `/page?utm_source=twitter` | `/page` | クエリパラメータ除去 |
| `/page?a=1&b=2` | `/page` | すべてのクエリパラメータ除去 |
| `/docs/guide/` | `/docs/guide` | トレイリングスラッシュ除去 |
| `/` | `/` | ルートパスは保持 |
| `/docs//guide` | `/docs/guide` | 連続スラッシュ正規化 |
| `/docs/../api` | `/api` | ドットセグメント解決 |

## 使用例

### 基本的な使用（変更なし）

```bash
# 既存のコマンドがそのまま使用可能
linked-url-lister https://handbook.gitlab.com/handbook/

# JSON出力
linked-url-lister --json https://docs.example.com/

# 深度制限
linked-url-lister --max-depth=2 https://docs.example.com/
```

### 正規化の効果を確認

```bash
# 詳細モードで正規化の効果を確認
linked-url-lister --verbose https://example.com/docs/ 2>&1 | grep "発見"
```

## トラブルシューティング

### 期待より少ないURLが出力される

正規化により重複が除去されているため、これは期待通りの動作です。正規化前の動作に戻す必要がある場合は、お問い合わせください。

### 特定のクエリパラメータを保持したい

現在のバージョンではすべてのクエリパラメータが除去されます。特定のパラメータを保持する必要がある場合は、機能リクエストとしてお知らせください。

## 関連ドキュメント

- [データモデル](./data-model.md) - 正規化処理フローの詳細
- [リサーチ](./research.md) - 実装の技術的決定
