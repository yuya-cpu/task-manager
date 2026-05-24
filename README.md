# task-manager

Gin と SQLite で作るタスク管理アプリです。  
JWT 認証付き REST API と、同梱の Web UI でタスクの作成・一覧・更新・削除ができます。

[![CI](https://github.com/yuya-cpu/task-manager/actions/workflows/ci.yml/badge.svg)](https://github.com/yuya-cpu/task-manager/actions/workflows/ci.yml)

## 機能一覧

| カテゴリ | 内容 |
|----------|------|
| タスク管理 | CRUD、期限・優先度・ステータス |
| 認証 | ユーザー登録 / ログイン（JWT）、ユーザーごとにタスクを分離 |
| 一覧 API | ステータス・優先度でのフィルタ、ソート、ページネーション |
| データ | SQLite 永続化（CGO 不要ドライバ） |
| フロント | 静的 Web UI（ログイン、タスク操作、絞り込み） |
| 品質 | 単体テスト・API 統合テスト |
| CI | GitHub Actions（`go test` / `go build` / Docker ビルド） |
| デプロイ | Docker / docker-compose 対応 |

## クイックスタート（Docker 推奨）

**必要:** [Docker Desktop](https://www.docker.com/products/docker-desktop/) を起動しておく

```bash
git clone https://github.com/yuya-cpu/task-manager.git
cd task-manager
docker compose up --build
```

| URL | 説明 |
|-----|------|
| http://127.0.0.1:8080/web/index.html | Web UI |
| http://127.0.0.1:8080 | API（ルートは Web UI へリダイレクト） |

### デモアカウント

| 項目 | 値 |
|------|-----|
| Email | `demo@example.com` |
| Password | `password123` |

初回起動時にデモユーザーとサンプルタスク 3 件が自動投入されます。

### 動作確認（Windows + Docker Desktop）

以下をローカルで確認済みです。

- `docker compose build` … イメージビルド成功
- `docker compose up -d` … コンテナ起動（ポート `8080`）
- `POST /auth/login` … ログイン成功
- `GET /assignments` … タスク一覧取得（`meta.total` 付き）
- `GET /web/index.html` … Web UI 表示（HTTP 200）

停止・削除:

```bash
docker compose down        # 停止
docker compose down -v     # 停止 + DB ボリューム削除
```

## ローカル開発（Go）

**必要:** Go 1.25 以降

```bash
cp .env.example .env   # SECRET_KEY を変更推奨
go mod download
go run .
```

- API / フロント: いずれも http://127.0.0.1:8080（Docker と同じ URL）

## 技術スタック

| 用途 | ライブラリ / ツール |
|------|---------------------|
| API | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io/) |
| DB | SQLite（[glebarez/sqlite](https://github.com/glebarez/sqlite)） |
| 認証 | [golang-jwt](https://github.com/golang-jwt/jwt) + bcrypt |
| テスト | testify + httptest |
| CI | GitHub Actions |
| コンテナ | Docker（multi-stage build） |
| フロント | HTML / CSS / JavaScript |

## アーキテクチャ

```
[Web UI]  web/index.html
    ↓  Authorization: Bearer <JWT>
[Gin API]  Handler → Service → Repository
    ↓
[SQLite]   data/task-manager.db（Docker 時はボリューム task-data）
```

## 認証（JWT）

1. `POST /auth/signup` でユーザー登録
2. `POST /auth/login` で `token` を取得
3. `/assignments` 系リクエストにヘッダを付与:

```
Authorization: Bearer <token>
```

## API リファレンス

### 認証（トークン不要）

| メソッド | パス | 説明 |
|----------|------|------|
| POST | `/auth/signup` | 新規登録 |
| POST | `/auth/login` | ログイン |

### タスク（要 JWT）

| メソッド | パス | 説明 |
|----------|------|------|
| GET | `/assignments` | 一覧（フィルタ・ページネーション） |
| GET | `/assignments/:id` | 1件取得 |
| POST | `/assignments` | 作成 |
| PUT | `/assignments/:id` | 更新 |
| DELETE | `/assignments/:id` | 削除 |

### 一覧クエリパラメータ

| パラメータ | 例 | 説明 |
|------------|-----|------|
| `status` | `todo` | `todo` / `in_progress` / `done` |
| `priority` | `high` | `low` / `medium` / `high` |
| `sort` | `newest` | `due_date_asc`（既定）/ `due_date_desc` / `newest` |
| `page` | `1` | ページ番号（既定: 1） |
| `limit` | `20` | 件数（既定: 20、最大: 100） |

レスポンス:

```json
{
  "data": [
    {
      "id": 1,
      "title": "Goの課題を提出する",
      "priority": "high",
      "status": "todo",
      "due_date": "2026-05-25T00:00:00Z"
    }
  ],
  "meta": { "page": 1, "limit": 20, "total": 3 }
}
```

### curl 例

```bash
# ログイン
curl -s -X POST http://127.0.0.1:8080/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"demo@example.com\",\"password\":\"password123\"}"

# 未完了タスクのみ（TOKEN を差し替え）
curl -s "http://127.0.0.1:8080/assignments?status=todo&sort=newest" \
  -H "Authorization: Bearer TOKEN"
```

## テスト

```bash
go test ./...
```

GitHub Actions（`.github/workflows/ci.yml`）では次を実行します。

- `go test ./...`
- `go build`
- `docker build`（イメージがビルドできることの確認）

## Docker 詳細

| ファイル | 役割 |
|----------|------|
| `Dockerfile` | Go マルチステージビルド → Alpine 実行イメージ |
| `docker-compose.yml` | ポート公開・環境変数・DB ボリューム |
| `.dockerignore` | ビルドコンテキストの除外 |

環境変数の上書き例:

```bash
# PowerShell
$env:SECRET_KEY="your-production-secret"; docker compose up --build

# bash
SECRET_KEY=your-production-secret docker compose up --build
```

## プロジェクト構成

```
task-manager/
├── .github/workflows/ci.yml   # CI
├── Dockerfile
├── docker-compose.yml
├── main.go / main_test.go
├── data/                      # DB 接続
├── models/                    # User, Assignment
├── dto/
├── handlers/
├── services/
├── repositories/
├── middlewares/               # JWT
└── web/                       # フロント
```

## 環境変数

| 変数 | デフォルト | 説明 |
|------|------------|------|
| `SECRET_KEY` | 開発用フォールバックあり | JWT 署名キー（本番では必ず変更） |
| `DB_PATH` | `data/task-manager.db` | SQLite ファイルパス |
| `GIN_MODE` | - | Docker 時は `release` |

## トラブルシューティング

### Docker が起動しない

```
failed to connect to the docker API ...
```

→ **Docker Desktop を起動**してから `docker compose up --build` を再実行してください。

### ポート 8080 が使用中

→ `docker-compose.yml` の `ports` を `"8081:8080"` などに変更するか、既存プロセスを停止してください。

### DB スキーマエラー（ローカル `go run` 時）

古い SQLite ファイルが原因のことがあります。

```bash
# Windows
del data\task-manager.db
go run .
```

Docker の DB をリセットする場合:

```bash
docker compose down -v
docker compose up --build
```

## 開発の経緯（参考）

1. Gin + SQLite でタスク CRUD API
2. JWT 認証・ユーザー別タスク分離
3. テスト・Web フロント追加
4. GitHub 公開・CI・クエリフィルタ・Docker 対応

## ライセンス

MIT
