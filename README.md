# task-manager

Gin と SQLite で作るタスク管理アプリです。  
REST API に **JWT 認証** を加え、同梱の **Web フロント** からタスクを操作できます。

## 機能

- ユーザー登録・ログイン（JWT）
- タスク（Assignment）の CRUD（ログインユーザーごとに分離）
- 期限・優先度・ステータス
- SQLite 永続化
- 単体テスト・API 統合テスト
- 簡易 Web UI（HTML + JavaScript）

## 技術スタック

| 用途 | ライブラリ |
|------|------------|
| API | Gin |
| ORM | GORM |
| DB | SQLite（glebarez/sqlite・CGO 不要） |
| 認証 | JWT（golang-jwt） |
| テスト | testify + httptest |
| フロント | 静的 HTML / CSS / JS |

## 必要環境

- Go 1.25 以降

## セットアップ

```bash
git clone https://github.com/yuya-cpu/task-manager.git
cd task-manager
cp .env.example .env   # SECRET_KEY を本番用に変更
go mod download
go run .
```

- API: http://127.0.0.1:8080  
- フロント: http://127.0.0.1:8080/web/index.html  

### デモアカウント

| 項目 | 値 |
|------|-----|
| Email | `demo@example.com` |
| Password | `password123` |

## テスト

```bash
go test ./...
```

## API

### 認証（トークン不要）

| メソッド | パス | 説明 |
|----------|------|------|
| POST | `/auth/signup` | 新規登録 |
| POST | `/auth/login` | ログイン（`token` を返す） |

### タスク（要 `Authorization: Bearer <token>`）

| メソッド | パス | 説明 |
|----------|------|------|
| GET | `/assignments` | 一覧 |
| GET | `/assignments/:id` | 1件取得 |
| POST | `/assignments` | 作成 |
| PUT | `/assignments/:id` | 更新 |
| DELETE | `/assignments/:id` | 削除 |

### タスクフィールド

| フィールド | 説明 |
|------------|------|
| `title` | タイトル（必須） |
| `description` | 説明 |
| `due_date` | `YYYY-MM-DD` |
| `priority` | `low` / `medium` / `high` |
| `status` | `todo` / `in_progress` / `done` |

### curl 例

```bash
# ログイン
curl -X POST http://127.0.0.1:8080/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"demo@example.com\",\"password\":\"password123\"}"

# 一覧（TOKEN を差し替え）
curl http://127.0.0.1:8080/assignments \
  -H "Authorization: Bearer TOKEN"
```

## プロジェクト構成

```
task-manager/
├── main.go
├── main_test.go
├── data/
├── models/          # User, Assignment
├── dto/
├── handlers/        # Auth, Assignment
├── services/
├── repositories/
├── middlewares/     # JWT 認証
└── web/             # フロント（静的ファイル）
```

## 環境変数

| 変数 | デフォルト | 説明 |
|------|------------|------|
| `SECRET_KEY` | （開発用フォールバックあり） | JWT 署名キー |
| `DB_PATH` | `data/task-manager.db` | SQLite パス |

## 注意

スキーマ変更後に起動エラーになる場合は、古い DB ファイルを削除して再起動してください。

```bash
rm data/task-manager.db   # Windows: del data\task-manager.db
go run .
```

## ライセンス

MIT
