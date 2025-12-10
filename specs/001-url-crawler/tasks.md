# タスク: NotebookLMインポート用URLクローラー

**入力**: `/specs/001-url-crawler/` の設計ドキュメント
**前提**: plan.md（必須）、spec.md（必須）、research.md、data-model.md、contracts/

**テスト**: 仕様書でテストは任意（オプション）と定義されているため、テストタスクは含めない

**構成**: タスクはユーザーストーリーごとにグループ化され、独立した実装・テストが可能

## フォーマット: `[ID] [P?] [Story] 説明`

- **[P]**: 並列実行可能（異なるファイル、依存関係なし）
- **[Story]**: 所属するユーザーストーリー（US1, US2, US3）
- 説明にはファイルパスを含める

## パス規約

- **Go プロジェクト**: `cmd/`（メインパッケージ）、`internal/`（内部パッケージ）

---

## Phase 1: セットアップ（共有インフラ）

**目的**: プロジェクト初期化と基本構造

- [x] T001 Goモジュールを初期化（go mod init github.com/gimKondo/linked-url-lister）
- [x] T002 [P] プロジェクトディレクトリ構造を作成（cmd/linked-url-lister/, internal/crawler/, internal/filter/, internal/output/, internal/robots/）
- [x] T003 [P] .gitignoreを作成（バイナリ、一時ファイル除外）

---

## Phase 2: 基盤（ブロッキング前提条件）

**目的**: すべてのユーザーストーリー実装前に完了必須のコアインフラ

**⚠️ 重要**: このフェーズが完了するまでユーザーストーリー作業は開始不可

- [x] T004 CrawlConfig構造体を作成（internal/crawler/config.go）- 開始URL、最大深度、最小テキスト長、遅延、robots.txt遵守フラグ、JSON出力、最大ページ数、User-Agent
- [x] T005 [P] CrawlURL構造体を作成（internal/crawler/url.go）- URL、深度、発見元URL
- [x] T006 [P] Page構造体を作成（internal/crawler/page.go）- URL、ステータスコード、テキスト長、リンク、エラー
- [x] T007 [P] CrawlResult構造体を作成（internal/crawler/result.go）- URL一覧、訪問数、フィルタ数、エラー、所要時間
- [x] T008 [P] CrawlError構造体を作成（internal/crawler/error.go）- URL、エラーメッセージ、ステータスコード
- [x] T009 RobotsTxt構造体と解析関数を作成（internal/robots/robots.go）- Disallow/Allowパス、CrawlDelay

**チェックポイント**: 基盤完了 - ユーザーストーリー実装を開始可能

---

## Phase 3: ユーザーストーリー1 - 基本的なURLクロール (優先度: P1) 🎯 MVP

**目標**: 開始URLから配下ページを再帰的にクロールし、テキストコンテンツでフィルタリングしたURLリストを出力

**独立テスト**: `linked-url-lister https://example.com/docs/` を実行し、配下URLがリストアップされることを確認

### ユーザーストーリー1の実装

- [x] T010 [US1] HTTPクライアント作成関数を実装（internal/crawler/http.go）- タイムアウト、User-Agent設定、リダイレクト追跡
- [x] T011 [US1] HTMLパーサーを実装（internal/crawler/parser.go）- golang.org/x/net/htmlを使用してリンク抽出
- [x] T012 [US1] テキスト抽出関数を実装（internal/filter/text.go）- script/style/noscriptタグ除外、空白正規化、文字数カウント
- [x] T013 [US1] URL正規化関数を実装（internal/filter/url.go）- 相対URL解決、フラグメント除去
- [x] T014 [US1] パスプレフィックスフィルタを実装（internal/filter/path.go）- 開始URLと同一ホスト・パス配下判定
- [x] T015 [US1] 訪問済みURL追跡を実装（internal/crawler/visited.go）- mapによる重複検出
- [x] T016 [US1] クローラーメインループを実装（internal/crawler/crawler.go）- キュー処理、フェッチ、フィルタ、リンク追加
- [x] T017 [US1] 進捗出力を実装（internal/crawler/progress.go）- stderrへの進捗メッセージ出力
- [x] T018 [US1] テキスト出力を実装（internal/output/text.go）- stdoutへの1行1URL出力
- [x] T019 [US1] CLIエントリーポイントを作成（cmd/linked-url-lister/main.go）- 引数パース、クローラー実行、結果出力
- [x] T020 [US1] 基本フラグを実装（cmd/linked-url-lister/main.go）- URL引数、--min-text、--help

**チェックポイント**: ユーザーストーリー1が独立して動作・テスト可能

---

## Phase 4: ユーザーストーリー2 - NotebookLM用の出力 (優先度: P2)

**目標**: NotebookLMに直接インポート可能な出力形式（テキスト/JSON）をサポート

**独立テスト**: `linked-url-lister --json https://example.com/` を実行し、JSON形式で出力されることを確認

### ユーザーストーリー2の実装

- [x] T021 [US2] JSON出力構造体を作成（internal/output/json.go）- urls配列、stats、errorsを含むJSONスキーマ
- [x] T022 [US2] JSON出力関数を実装（internal/output/json.go）- encoding/jsonでマーシャル、stdoutへ出力
- [x] T023 [US2] --jsonフラグを追加（cmd/linked-url-lister/main.go）- OutputJSONフラグ処理
- [x] T024 [US2] 統計情報収集を実装（internal/crawler/crawler.go）- 訪問数、フィルタ数、所要時間をCrawlResultに追加

**チェックポイント**: ユーザーストーリー1と2が独立して動作

---

## Phase 5: ユーザーストーリー3 - クロール制御 (優先度: P3)

**目標**: 最大深度、リクエスト間隔、robots.txt遵守などのクロール制御オプション

**独立テスト**: `linked-url-lister --max-depth=2 --delay=500ms https://example.com/` で動作が変化することを確認

### ユーザーストーリー3の実装

- [x] T025 [US3] 深度制限を実装（internal/crawler/crawler.go）- MaxDepth超過時のスキップ
- [x] T026 [US3] リクエスト遅延を実装（internal/crawler/crawler.go）- time.Sleep による間隔制御
- [x] T027 [US3] robots.txt取得・解析を実装（internal/robots/fetch.go）- ホストごとにキャッシュ
- [x] T028 [US3] robots.txtチェックを統合（internal/crawler/crawler.go）- クロール前にパス許可判定
- [x] T029 [US3] 最大ページ数制限を実装（internal/crawler/crawler.go）- MaxPages超過時の終了
- [x] T030 [US3] 追加フラグを実装（cmd/linked-url-lister/main.go）- --max-depth, --delay, --max-pages, --ignore-robots, --user-agent, --verbose
- [x] T031 [US3] 詳細進捗出力を実装（internal/crawler/progress.go）- --verbose時のリンク数・テキスト長表示

**チェックポイント**: すべてのユーザーストーリーが独立して動作

---

## Phase 6: 仕上げと横断的関心事

**目的**: 複数のユーザーストーリーに影響する改善

- [x] T032 [P] READMEを作成（README.md）- 目的、インストール、使用例、貢献ガイドライン
- [x] T033 [P] --helpテキストを充実（cmd/linked-url-lister/main.go）- 全オプションの説明
- [x] T034 [P] --versionフラグを実装（cmd/linked-url-lister/main.go）- バージョン表示
- [x] T035 エラーハンドリングを強化（全internal/パッケージ）- ユーザーフレンドリーなエラーメッセージ
- [x] T036 終了コードを実装（cmd/linked-url-lister/main.go）- 0:成功、1:引数エラー、2:ネットワークエラー、3:中断
- [x] T037 quickstart.md検証を実行 - ドキュメントの手順が動作することを確認

---

## 依存関係と実行順序

### フェーズ依存関係

- **セットアップ (Phase 1)**: 依存なし - 即時開始可能
- **基盤 (Phase 2)**: セットアップ完了に依存 - すべてのユーザーストーリーをブロック
- **ユーザーストーリー (Phase 3+)**: 基盤フェーズ完了に依存
  - ストーリーは優先度順に実行（P1 → P2 → P3）
  - または並列実行可能（リソースがあれば）
- **仕上げ (Phase 6)**: 必要なユーザーストーリー完了に依存

### ユーザーストーリー依存関係

- **ユーザーストーリー1 (P1)**: 基盤完了後に開始可能 - 他ストーリーへの依存なし
- **ユーザーストーリー2 (P2)**: US1の出力インフラを使用するが独立テスト可能
- **ユーザーストーリー3 (P3)**: US1のクローラーコアを拡張するが独立テスト可能

### 各ユーザーストーリー内

- モデル → サービス → 出力の順
- コア実装 → 統合の順
- ストーリー完了後に次の優先度へ

### 並列実行機会

- Phase 1: T002, T003 は並列可能
- Phase 2: T005, T006, T007, T008 は並列可能
- Phase 6: T032, T033, T034 は並列可能

---

## 並列実行例: Phase 2

```bash
# Phase 2の構造体作成を並列実行:
Task: "CrawlURL構造体を作成 (internal/crawler/url.go)"
Task: "Page構造体を作成 (internal/crawler/page.go)"
Task: "CrawlResult構造体を作成 (internal/crawler/result.go)"
Task: "CrawlError構造体を作成 (internal/crawler/error.go)"
```

---

## 実装戦略

### MVP優先 (ユーザーストーリー1のみ)

1. Phase 1: セットアップ完了
2. Phase 2: 基盤完了（重要 - すべてのストーリーをブロック）
3. Phase 3: ユーザーストーリー1完了
4. **停止して検証**: ユーザーストーリー1を独立テスト
5. 準備できたらリリース/デモ

### インクリメンタルデリバリー

1. セットアップ + 基盤完了 → 基盤準備完了
2. ユーザーストーリー1追加 → 独立テスト → リリース/デモ (MVP!)
3. ユーザーストーリー2追加 → 独立テスト → リリース/デモ
4. ユーザーストーリー3追加 → 独立テスト → リリース/デモ
5. 各ストーリーは前のストーリーを壊さずに価値を追加

---

## 備考

- [P] タスク = 異なるファイル、依存関係なし
- [Story] ラベル = 特定のユーザーストーリーへの追跡性
- 各ユーザーストーリーは独立して完了・テスト可能であるべき
- タスクまたは論理グループごとにコミット
- チェックポイントで停止してストーリーを独立検証
- 回避すべき: 曖昧なタスク、同一ファイル競合、独立性を損なうストーリー間依存
