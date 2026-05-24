# task-manager

Go（Gin）と SQLite で作った**タスク管理 Web アプリ**です。  
ブラウザからタスクの登録・一覧・更新・削除ができ、バックエンドは REST API として動作します。

[![CI](https://github.com/yuya-cpu/task-manager/actions/workflows/ci.yml/badge.svg)](https://github.com/yuya-cpu/task-manager/actions/workflows/ci.yml)

## このプロジェクトについて

個人のタスク（課題）を管理するためのアプリです。  
1 ユーザーごとにタスクが分かれており、ログインしないと他人のタスクは見えません。

**主な機能**

- タスクの作成・一覧・更新・削除
- 期限（`due_date`）、優先度（`low` / `medium` / `high`）、ステータス（`todo` / `in_progress` / `done`）
- ユーザー登録・ログイン（JWT 認証）
- 一覧の絞り込み（ステータス・優先度）と並び替え
- Web UI（`web/`）から操作可能

**技術的な特徴**

- API は **Handler → Service → Repository** の 3 層構成
- 単体テスト・API テストあり
- GitHub Actions で CI（テスト・ビルド・Docker ビルド）
- `docker compose up` で起動可能

## 使い方

### 1. 起動（Docker が簡単）

[Docker Desktop](https://www.docker.com/products/docker-desktop/) を起動したうえで:

```bash
git clone https://github.com/yuya-cpu/task-manager.git
cd task-manager
docker compose up --build
```

ブラウザで開く: **http://127.0.0.1:8080/web/index.html**

### 2. ログイン

初回起動時にデモ用のユーザーとサンプルタスクが入ります。

| Email | Password |
|-------|----------|
| `demo@example.com` | `password123` |

新規登録も Web UI または API から可能です。

### 3. ローカルで起動（Go）

Go 1.25 以降が必要です。

```bash
go mod download
go run .
```

API と Web UI はどちらも `http://127.0.0.1:8080` で利用できます。

## 技術スタック

| 区分 | 技術 |
|------|------|
| 言語 | Go |
| Web フレームワーク | Gin |
| DB | SQLite（GORM） |
| 認証 | JWT + bcrypt |
| フロント | HTML / CSS / JavaScript（ビルド不要） |
| テスト | testify, httptest |
| CI | GitHub Actions |
| コンテナ | Docker |

## ディレクトリ構成

```
task-manager/
├── main.go              # エントリポイント・ルーティング
├── handlers/            # HTTP リクエストの受付
├── services/            # ビジネスロジック
├── repositories/        # DB 操作
├── models/              # データ構造
├── middlewares/         # JWT 認証
├── web/                 # フロント（静的ファイル）
├── Dockerfile
└── docker-compose.yml
```

## API 概要

認証後、リクエストヘッダにトークンを付けます。

```
Authorization: Bearer <ログインで取得した token>
```

### 認証（トークン不要）

| メソッド | パス | 説明 |
|----------|------|------|
| POST | `/auth/signup` | 新規登録 |
| POST | `/auth/login` | ログイン（`token` を返す） |

### タスク（トークン必須）

| メソッド | パス | 説明 |
|----------|------|------|
| GET | `/assignments` | 一覧 |
| GET | `/assignments/:id` | 1 件取得 |
| POST | `/assignments` | 作成 |
| PUT | `/assignments/:id` | 更新 |
| DELETE | `/assignments/:id` | 削除 |

**一覧のクエリ例**

```
GET /assignments?status=todo&priority=high&sort=newest&page=1&limit=20
```

レスポンスには `data`（タスク配列）と `meta`（件数・ページ情報）が含まれます。

```json
{
  "data": [
    {
      "id": 1,
      "title": "レポート提出",
      "priority": "high",
      "status": "todo",
      "due_date": "2026-06-01T00:00:00Z"
    }
  ],
  "meta": { "page": 1, "limit": 20, "total": 1 }
}
```

## テスト

```bash
go test ./...
```

`main` ブランチへの push 時に GitHub Actions でも同様のチェックが走ります。

## ライセンス

MIT
