# syntax=docker/dockerfile:1.19
# 開発用イメージ。本番用の multi-stage（standalone 出力）は、デプロイ先が決まってから足す。
FROM node:24.22.1-alpine3.24

ENV NEXT_TELEMETRY_DISABLED=1

WORKDIR /app

# node 公式イメージには uid 1000 の node ユーザーが既にいるので、それを使う。
# 依存のインストールより前に USER を切り替える。root で npm ci すると
# node_modules が root 所有で作られ、
#   - コンテナ内から npm install できない
#   - 依存の install スクリプトが root で走る
# の2つを踏む。
RUN chown node:node /app
USER node

# 依存を先に入れてレイヤーキャッシュを効かせる。
# ここで入れた node_modules が、compose の named volume の初期値になる。
COPY --chown=node:node package.json package-lock.json ./
RUN npm ci

# .next も named volume にするため、先に作っておく（node 所有になる）。
RUN mkdir -p /app/.next

EXPOSE 3000
CMD ["npm", "run", "dev"]
