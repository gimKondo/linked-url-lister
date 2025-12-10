# 実装計画: NotebookLMインポート用URLクローラー

**ブランチ**: `001-url-crawler` | **日付**: 2025-12-10 | **仕様書**: [spec.md](./spec.md)
**入力**: `/specs/001-url-crawler/spec.md` の機能仕様書

## 概要

指定したWebサイトの配下ページURLをリストアップするCLIツールを実装する。主な用途はNotebookLMへのURL一括インポート。開始URLから再帰的にリンクを辿り、同一パス配下かつ十分なテキストコンテンツを持つページのURLを抽出・出力する。

## 技術コンテキスト

**言語/バージョン**: Go 1.21以上（標準ライブラリの充実度を考慮）
**主要依存関係**: 標準ライブラリのみ（net/http, html, encoding/json）、外部依存なし
**ストレージ**: N/A（ステートレスなCLIツール）
**テスト**: go test（オプション、憲法に準拠）
**対象プラットフォーム**: クロスプラットフォーム（Linux, macOS, Windows）
**プロジェクトタイプ**: シングルプロジェクト
**パフォーマンス目標**: 100ページを5分以内に処理（SC-001より）
**制約**: リクエスト間隔デフォルト100ms、メモリ使用量は訪問済みURL数に比例
**スケール/スコープ**: 数百〜数千ページのドキュメントサイト

## 憲法チェック

*ゲート: Phase 0リサーチ前に通過必須。Phase 1設計後に再チェック。*

| 原則 | 状態 | 検証 |
|------|------|------|
| I. シンプルさ優先 | ✅ 合格 | 外部依存なし、標準ライブラリのみ使用。単一バイナリで完結。 |
| II. CLI優先インターフェース | ✅ 合格 | stdin/stdout/stderr分離、JSON対応、POSIX準拠フラグ |
| III. ドキュメント | ✅ 合格 | --help必須、README作成予定 |

## プロジェクト構造

### ドキュメント（この機能）

```text
specs/001-url-crawler/
├── plan.md              # このファイル
├── research.md          # Phase 0出力
├── data-model.md        # Phase 1出力
├── quickstart.md        # Phase 1出力
├── contracts/           # Phase 1出力（CLIインターフェース定義）
└── tasks.md             # Phase 2出力（/speckit.tasksで作成）
```

### ソースコード（リポジトリルート）

```text
cmd/
└── linked-url-lister/
    └── main.go          # エントリーポイント

internal/
├── crawler/
│   └── crawler.go       # クロールロジック
├── filter/
│   └── filter.go        # URL/コンテンツフィルタリング
├── output/
│   └── output.go        # 出力フォーマット（テキスト/JSON）
└── robots/
    └── robots.go        # robots.txt解析
```

**構造決定**: Goの標準的なプロジェクトレイアウトを採用。`cmd/`にメインパッケージ、`internal/`に内部パッケージを配置。パッケージ数は最小限（4パッケージ）に抑え、シンプルさを維持。

## 複雑性追跡

> 憲法違反がない場合は空欄

該当なし - すべての設計選択は憲法に準拠している。
