# warican

Goとなんかで作る割り勘アプリ

## 構成

```
warican/
├── compose.yaml          # web / api の2サービス
├── .env.example          # ポートと uid/gid。cp して .env にする
├── api/                  # Go (net/http)。air でホットリロード
│   ├── Dockerfile
│   ├── .air.toml
│   ├── go.mod
│   └── cmd/server/main.go
└── web/                  # Next.js (App Router, TypeScript)
    ├── Dockerfile
    └── app/
```

ホストに必要なのは Docker だけ。Go と Node はホストに入れない。

## 起動

```sh
cp .env.example .env     # 任意。無くても既定値で動く
docker compose up --build
```

| | URL |
|---|---|
| web | http://localhost:3000 |
| api | http://localhost:8080/healthz |

`http://localhost:3000` を開くと、**web コンテナ → api** と **ブラウザ → api** の
2経路の疎通結果が出る。両方 `200` なら環境は通っている。
このページ（`web/app/page.tsx` と `web/app/health-check.tsx`）はアプリを書き始めたら消してよい。

停止は `docker compose down`。volume ごと消すなら `docker compose down -v`。

## バージョン（すべて固定。`latest` は使わない）

| 対象 | バージョン | 選定理由 |
|---|---|---|
| Go | `1.27.3-alpine3.24` | 2026-10-19 リリースの最新安定 |
| Node | `24.22.1-alpine3.24` | 24系 (Krypton) が最新 LTS。26系は非LTS |
| Next.js | `16.4.3` | npm の最新安定 |
| React | `19.3.0` | Next 16.4.3 のスキャフォールドが固定した版 |
| air | `v1.67.4` | Go ホットリロード。`api/Dockerfile` の `AIR_VERSION` |

`web/package.json` の devDependencies だけ `^` 付きだが、`package-lock.json` が
解決後の版を固定しているので `npm ci` で再現する。

## 詰まりやすいところ

### API の URL が 2 本ある

| 変数 | 値 | 誰が使うか |
|---|---|---|
| `API_URL` | `http://api:8080` | web コンテナの中（SSR / Server Components） |
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080` | ブラウザ（クライアントコンポーネント） |

ブラウザから見た `localhost` はホスト、web コンテナから見た `localhost` は
web コンテナ自身。同じ値は使えない。`NEXT_PUBLIC_` 付きの値はクライアント JS に
文字列として埋め込まれるので、秘密は入れない。

ブラウザ → api はクロスオリジンになるので、api 側に CORS が必要。
許可オリジンは `CORS_ALLOWED_ORIGIN` の1本だけ（`compose.yaml` で `WEB_PORT` と連動）。

### uid/gid を揃えている

`api` コンテナは uid 1000 で動かしている。root で動かすと bind mount 上に
root 所有のファイル（`api/tmp/`, `go.sum` など）ができてホストから消せなくなる。
ホストの uid が 1000 でない場合は `.env` の `UID` / `GID` を `id -u` / `id -g` に合わせる。

### node_modules と .next は named volume

`./web` を丸ごと bind mount すると、ホスト側の（存在しない）`node_modules` で
コンテナ内のそれが上書きされる。`compose.yaml` で両ディレクトリに named volume を
重ねて回避している。初期値は `web/Dockerfile` の `npm ci` が作ったものが入る。

そのため **`web/package.json` を変えたときは volume の作り直しが必要**:

```sh
docker compose down
docker volume rm warican_web-node-modules
docker compose up --build
```

### ホットリロードが効かなくなったら

`web` は Turbopack の**内容ハッシュ**で変更を判定する。`touch` だけ（内容が同じ）では
再コンパイルされないのは正常。

長時間動かしていると、ごくまれに編集を取りこぼしてブラウザが古い内容のまま固まることがある。
実測では、コンテナ再起動直後は32回連続の編集で取りこぼしゼロ、
長時間動かしたコンテナで16回中1回（次の編集で自動的に復帰）。固まったら:

```sh
docker compose restart web
```

### ポートを変えるとき

`.env` の `WEB_PORT` / `API_PORT` だけ変える。`CORS_ALLOWED_ORIGIN` と
`NEXT_PUBLIC_API_URL` は `compose.yaml` 側で連動している。

## まだ決めていないこと

| 項目 | 現状 |
|---|---|
| DB | 未導入。MySQL 9.7.4 (LTS) 予定 → 下記手順 |
| マイグレーションツール | 未選定（golang-migrate / goose / Atlas） |
| Go の HTTP 層 | 標準 `net/http` の `ServeMux` のみ。外部依存ゼロ（`go.mod` 参照） |
| CSS | Tailwind は入れていない。素の CSS |
| 本番用ビルド | 未作成。デプロイ先が決まってから multi-stage を足す |

### MySQL を足すとき

`compose.yaml` に追記する:

```yaml
  db:
    image: mysql:9.7.4
    environment:
      MYSQL_DATABASE: warican
      MYSQL_USER: warican
      MYSQL_PASSWORD: ${DB_PASSWORD:-warican}
      MYSQL_ROOT_PASSWORD: ${DB_ROOT_PASSWORD:-root}
    ports:
      - "${DB_PORT:-3306}:3306"
    volumes:
      - mysql-data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "127.0.0.1", "-p${MYSQL_ROOT_PASSWORD}"]
      interval: 10s
      timeout: 5s
      retries: 10
      start_period: 60s
```

同時に必要になるのは3つ:
- `volumes:` に `mysql-data:` を追加
- `api` に `depends_on: db: condition: service_healthy` と DSN の環境変数
- Go 側のドライバ（`go-sql-driver/mysql`）と、マイグレーションの流し方
