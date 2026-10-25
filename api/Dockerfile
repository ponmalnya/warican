# syntax=docker/dockerfile:1.19
# 開発用イメージ。本番用の multi-stage は、デプロイ先が決まってから足す。
FROM golang:1.27.3-alpine3.24

# ホスト（Linux）と uid を揃える。
# 揃えないと、bind mount 上に root 所有のファイル（tmp/, go.sum 等）が
# できてホスト側から消せなくなる。
ARG UID=1000
ARG GID=1000
RUN addgroup -g ${GID} app \
 && adduser -u ${UID} -G app -s /bin/sh -D app \
 && mkdir -p /go/pkg/mod /home/app/.cache/go-build \
 && chown -R ${UID}:${GID} /go /home/app

ENV GOCACHE=/home/app/.cache/go-build \
    CGO_ENABLED=0

# 依存のビルドより前に USER を切り替える。サードパーティのコードを
# root でコンパイルしないため。
USER app

# air（Go のホットリロード）。バージョンを固定して再現性を確保する。
ARG AIR_VERSION=v1.67.4
RUN go install github.com/air-verse/air@${AIR_VERSION}

WORKDIR /app
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]
