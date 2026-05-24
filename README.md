# task-manager

Gin と SQLite で作る、シンプルなタスク管理 REST API です。  
タスク（Assignment）の作成・一覧・更新・削除ができ、**期限・優先度・ステータス** を扱えます。

## 機能

- タスク（Assignment）の CRUD
- 期限（`due_date`）・優先度（`priority`）・ステータス（`status`）
- SQLite への永続化（ファイル: `data/task-manager.db`）
- 初回起動時のサンプルデータ投入

## 技術スタック

| 用途 | ライブラリ |
|------|------------|
| Web フレームワーク | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io/) |
| DB | SQLite（[glebarez/sqlite](https://github.com/glebarez/sqlite)・CGO 不要） |

## 必要環境

- Go 1.25 以降

## セットアップ

```bash
git clone https://github.com/yuya-cpu/task-manager.git
cd task-manager
go mod download
go run .
```

起動後: `http://127.0.0.1:8080`

## API

ベース URL: `http://127.0.0.1:8080`

| メソッド | パス | 説明 |
|----------|------|------|
| GET | `/assignments` | 一覧取得 |
| GET | `/assignments/:id` | 1件取得 |
| POST | `/assignments` | 作成 |
| PUT | `/assignments/:id` | 更新 |
| DELETE | `/assignments/:id` | 削除 |

### フィールド

| フィールド | 型 | 説明 |
|------------|-----|------|
| `title` | string | タイトル（必須） |
| `description` | string | 説明 |
| `due_date` | string | 期限（`YYYY-MM-DD`、省略可） |
| `priority` | string | `low` / `medium` / `high` |
| `status` | string | `todo` / `in_progress` / `done` |

### 作成例

```bash
curl -X POST http://127.0.0.1:8080/assignments \
  -H "Content-Type: application/json" \
  -d "{\"title\":\"レポート提出\",\"description\":\"第3章まで\",\"due_date\":\"2026-06-01\",\"priority\":\"high\",\"status\":\"todo\"}"
```

### 更新例

```bash
curl -X PUT http://127.0.0.1:8080/assignments/1 \
  -H "Content-Type: application/json" \
  -d "{\"status\":\"done\"}"
```

## プロジェクト構成

```
task-manager/
├── main.go                 # 起動・ルーティング・シード
├── data/db.go              # DB 接続・マイグレーション
├── models/assignment.go    # データモデル
├── dto/assignment_dto.go   # リクエスト用 DTO
├── handlers/assignment.go  # HTTP ハンドラ
├── services/               # ビジネスロジック
└── repositories/           # DB アクセス
```

**Handler → Service → Repository** の3層構成です（`gin-fleamarket` と同じ考え方）。

## 環境変数

| 変数 | デフォルト | 説明 |
|------|------------|------|
| `DB_PATH` | `data/task-manager.db` | SQLite ファイルのパス |

## ライセンス

MIT（必要に応じて変更してください）
