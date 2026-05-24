# task-manager

Gin + SQLite のタスク管理 API。JWT 認証、Web UI、テスト、CI、Docker 対応。

[![CI](https://github.com/yuya-cpu/task-manager/actions/workflows/ci.yml/badge.svg)](https://github.com/yuya-cpu/task-manager/actions/workflows/ci.yml)

## できること

- タスクの CRUD（期限・優先度・ステータス）
- ユーザー登録 / ログイン（JWT）
- 一覧のフィルタ・ソート・ページネーション
- Web UI から操作可能

## 起動

**Docker（推奨）**

```bash
git clone https://github.com/yuya-cpu/task-manager.git
cd task-manager
docker compose up --build
```

→ http://127.0.0.1:8080/web/index.html

**ローカル**

```bash
go run .
```

### デモアカウント

| Email | Password |
|-------|----------|
| `demo@example.com` | `password123` |

## 技術スタック

Go / Gin / GORM / SQLite / JWT / testify / GitHub Actions / Docker

## 構成

```
Handler → Service → Repository → SQLite
```

## API（概要）

| メソッド | パス | 認証 |
|----------|------|------|
| POST | `/auth/signup`, `/auth/login` | 不要 |
| GET/POST | `/assignments` | JWT 必須 |
| GET/PUT/DELETE | `/assignments/:id` | JWT 必須 |

一覧例: `GET /assignments?status=todo&priority=high&sort=newest&page=1&limit=20`

## テスト

```bash
go test ./...
```

## ライセンス

MIT
